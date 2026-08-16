// Package githooks holds nothing but the test for this repository's bootstrap git hooks, so
// `go test ./...` covers them alongside the SDK.
//
// .githooks/pre-commit and .githooks/post-checkout are the only hooks this repository authors;
// the rest of .githooks/ and the .ai-shared/ cache they fetch into are written by provisioning
// from that cache. Each wrapper is thin, so the URL it names, the directory it caches into, what
// it leaves behind when a transfer dies, and what it says when the fetch fails are the whole of
// its behavior, and all four are asserted here. What this file pins for the generated hooks is
// that .gitignore keeps them out.
package githooks

import (
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	publishingRepo              = "mvdatacenter/ai-instructions"
	cacheDir                    = ".ai-shared/"
	repoTheRenameLeftBehind     = "mvdatacenter/claude-instructions"
	cacheDirTheRenameLeftBehind = ".claude-shared/"
)

// curlDiesMidTransfer stands in for a transfer that dies after the server started sending.
// Measured against a real curl and a local HTTP server that declares Content-Length 5000, sends
// 51 bytes and closes: curl writes those 51 bytes to the -o path and exits 18, CURLE_PARTIAL_FILE.
// That leftover is what a wrapper fetching straight to the cache path keeps and runs for good,
// because a wrapper fetches only when the cache path is absent. An HTTP error is a different case
// and not the dangerous one: the same curl against a 404 under -f exits 22 and creates no file.
const curlDiesMidTransfer = `#!/bin/sh
out=""
while [ $# -gt 0 ]; do
    if [ "$1" = "-o" ]; then out="$2"; fi
    shift
done
[ -n "$out" ] && printf '#!/bin/bash\n# Central logic, the first 2 kB of it' > "$out"
exit 18
`

const ghAnswersNothing = "#!/bin/sh\nexit 0\n"

// The two wrappers this repository authors, and what each owes a fetch it could not complete.
var bootstrapWrappers = []struct {
	name         string
	cachedScript string
	// pre-commit refuses the commit, because the central staged-diff credential scan it needs
	// never ran. post-checkout reports and exits 0, because the checkout has already happened
	// and the hook's status is the status git gives back for the clone or switch.
	refusesTheOperation bool
}{
	{name: "pre-commit", cachedScript: cacheDir + "pre-commit.sh", refusesTheOperation: true},
	{name: "post-checkout", cachedScript: cacheDir + "post-checkout.sh", refusesTheOperation: false},
}

// Written by provisioning, never authored here. A tracked copy is whatever was last committed
// until provisioning overwrites it, which is how this repo carried pre-rename fetchers.
var generatedWrappers = []string{
	".githooks/commit-msg",
	".githooks/pre-push",
	".githooks/guard-destructive.sh",
	".githooks/session-lock.sh",
}

func wrapperPath(t *testing.T, name string) string {
	t.Helper()
	path, err := filepath.Abs(filepath.Join(".githooks", name))
	if err != nil {
		t.Fatalf("resolving .githooks/%s: %v", name, err)
	}
	return path
}

func writeStub(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatalf("writing stub %s: %v", path, err)
	}
}

// runAgainstADyingTransfer runs a wrapper in an empty directory whose curl always dies
// mid-transfer, and answers with the directory it ran in so a caller can ask what survived.
func runAgainstADyingTransfer(t *testing.T, wrapper string) (workdir string, exitCode int, stderr string) {
	t.Helper()
	root := t.TempDir()
	stubs := filepath.Join(root, "stubs")
	workdir = filepath.Join(root, "checkout")
	for _, dir := range []string{stubs, workdir} {
		if err := os.Mkdir(dir, 0o755); err != nil {
			t.Fatalf("creating %s: %v", dir, err)
		}
	}
	writeStub(t, filepath.Join(stubs, "curl"), curlDiesMidTransfer)
	writeStub(t, filepath.Join(stubs, "gh"), ghAnswersNothing)

	var said strings.Builder
	hook := exec.Command("bash", wrapperPath(t, wrapper))
	hook.Dir = workdir
	hook.Env = append(os.Environ(), "PATH="+stubs+string(os.PathListSeparator)+os.Getenv("PATH"))
	hook.Stdout = io.Discard
	hook.Stderr = &said
	if err := hook.Run(); err != nil {
		if _, isExit := err.(*exec.ExitError); !isExit {
			t.Fatalf("running .githooks/%s: %v", wrapper, err)
		}
	}
	return workdir, hook.ProcessState.ExitCode(), said.String()
}

func filesLeftIn(t *testing.T, dir string) []string {
	t.Helper()
	var kept []string
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			relative, relErr := filepath.Rel(dir, path)
			if relErr != nil {
				return relErr
			}
			kept = append(kept, relative)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("reading what %s holds: %v", dir, err)
	}
	return kept
}

func TestTheWrappersFetchFromTheCurrentInstructionRepo(t *testing.T) {
	for _, wrapper := range bootstrapWrappers {
		t.Run(wrapper.name, func(t *testing.T) {
			text := readWrapper(t, wrapper.name)
			if !strings.Contains(text, publishingRepo) {
				t.Errorf(".githooks/%s does not fetch from %s", wrapper.name, publishingRepo)
			}
			if !strings.Contains(text, cacheDir) {
				t.Errorf(".githooks/%s does not cache into %s", wrapper.name, cacheDir)
			}
		})
	}
}

func TestTheWrappersDoNotNameWhatTheRenameLeftBehind(t *testing.T) {
	// Asked separately from the assertion above because a wrapper can carry both: a fetch that
	// moved to the new repo while the cache path stayed on the old one still reads a stale
	// script on every box that has one, and never refreshes it.
	for _, wrapper := range bootstrapWrappers {
		t.Run(wrapper.name, func(t *testing.T) {
			text := readWrapper(t, wrapper.name)
			if strings.Contains(text, repoTheRenameLeftBehind) {
				t.Errorf(".githooks/%s still names %s", wrapper.name, repoTheRenameLeftBehind)
			}
			if strings.Contains(text, cacheDirTheRenameLeftBehind) {
				t.Errorf(".githooks/%s still caches into %s", wrapper.name, cacheDirTheRenameLeftBehind)
			}
		})
	}
}

func readWrapper(t *testing.T, name string) string {
	t.Helper()
	text, err := os.ReadFile(wrapperPath(t, name))
	if err != nil {
		t.Fatalf("reading .githooks/%s: %v", name, err)
	}
	return string(text)
}

func TestNoPartOfADeadTransferIsKept(t *testing.T) {
	// Asked of the whole directory rather than of the cached script, because the file that must
	// not survive is whatever path the wrapper chose to write: a wrapper that keeps a fragment
	// under some other name has the same defect and would pass a named check.
	for _, wrapper := range bootstrapWrappers {
		t.Run(wrapper.name, func(t *testing.T) {
			workdir, _, _ := runAgainstADyingTransfer(t, wrapper.name)
			kept := filesLeftIn(t, workdir)
			if len(kept) == 0 {
				return
			}
			body, err := os.ReadFile(filepath.Join(workdir, kept[0]))
			if err != nil {
				t.Fatalf("reading %s: %v", kept[0], err)
			}
			t.Errorf(".githooks/%s kept %s from a transfer that died, and nothing refetches a "+
				"cache path that is already there: %q", wrapper.name, kept[0], string(body))
		})
	}
}

func TestADeadTransferIsSaidRatherThanPassedOver(t *testing.T) {
	for _, wrapper := range bootstrapWrappers {
		t.Run(wrapper.name, func(t *testing.T) {
			_, exitCode, stderr := runAgainstADyingTransfer(t, wrapper.name)
			if wrapper.refusesTheOperation && exitCode == 0 {
				t.Errorf(".githooks/%s reported success though the central checks never ran", wrapper.name)
			}
			if !wrapper.refusesTheOperation && exitCode != 0 {
				t.Errorf(".githooks/%s exited %d, turning a completed checkout into a failed one",
					wrapper.name, exitCode)
			}
			assertSaysWhatItCouldNotReach(t, wrapper.name, wrapper.cachedScript, stderr)
		})
	}
}

func TestTheWrappersSayWhyWhenThereIsNoCurl(t *testing.T) {
	// A PATH holding everything the wrappers need except a network client, so the one thing
	// missing is the fetch itself and no other absent tool can produce the refusal.
	needed := []string{"bash", "mkdir", "mv", "rm"}
	binaries := t.TempDir()
	for _, name := range needed {
		found, err := exec.LookPath(name)
		if err != nil {
			t.Fatalf("%s is required to run the hook wrappers: %v", name, err)
		}
		if err := os.Symlink(found, filepath.Join(binaries, name)); err != nil {
			t.Fatalf("linking %s: %v", name, err)
		}
	}
	for _, wrapper := range bootstrapWrappers {
		t.Run(wrapper.name, func(t *testing.T) {
			workdir := t.TempDir()
			var said strings.Builder
			hook := exec.Command(filepath.Join(binaries, "bash"), wrapperPath(t, wrapper.name))
			hook.Dir = workdir
			hook.Env = []string{"PATH=" + binaries}
			hook.Stdout = io.Discard
			hook.Stderr = &said
			if err := hook.Run(); err != nil {
				if _, isExit := err.(*exec.ExitError); !isExit {
					t.Fatalf("running .githooks/%s: %v", wrapper.name, err)
				}
			}
			if wrapper.refusesTheOperation && hook.ProcessState.ExitCode() == 0 {
				t.Errorf(".githooks/%s reported success with no way to fetch anything", wrapper.name)
			}
			assertSaysWhatItCouldNotReach(t, wrapper.name, wrapper.cachedScript, said.String())
		})
	}
}

func assertSaysWhatItCouldNotReach(t *testing.T, wrapper, cachedScript, stderr string) {
	t.Helper()
	if !strings.Contains(stderr, cachedScript) {
		t.Errorf(".githooks/%s does not name %s in what it said: %q", wrapper, cachedScript, stderr)
	}
	if !strings.Contains(stderr, publishingRepo) {
		t.Errorf(".githooks/%s does not name the URL it could not reach: %q", wrapper, stderr)
	}
}

// ignoredIn answers, for each path, whether this repository's .gitignore ignores it. Every answer
// comes from a throwaway repository holding only that file, never from the checkout it lives in,
// so a local exclude file cannot make the tracked .gitignore look right.
func ignoredIn(t *testing.T, paths []string) map[string]bool {
	t.Helper()
	scratch := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", scratch).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	rules, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}
	if err := os.WriteFile(filepath.Join(scratch, ".gitignore"), rules, 0o644); err != nil {
		t.Fatalf("writing .gitignore into the scratch repository: %v", err)
	}
	answers := make(map[string]bool, len(paths))
	for _, path := range paths {
		target := filepath.Join(scratch, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatalf("creating %s: %v", path, err)
		}
		if err := os.WriteFile(target, nil, 0o644); err != nil {
			t.Fatalf("creating %s: %v", path, err)
		}
		asked := exec.Command("git", "check-ignore", "-q", path)
		asked.Dir = scratch
		answers[path] = asked.Run() == nil
	}
	return answers
}

func TestGitignoreKeepsTheCacheAndTheGeneratedWrappersOut(t *testing.T) {
	var asked []string
	for _, wrapper := range bootstrapWrappers {
		asked = append(asked, wrapper.cachedScript)
	}
	// The pre-rename cache path stays ignored so a checkout that still carries one is not
	// turned into an untracked-file failure by the rename.
	asked = append(asked, cacheDirTheRenameLeftBehind+"pre-commit.sh")
	asked = append(asked, generatedWrappers...)
	for path, ignored := range ignoredIn(t, asked) {
		if !ignored {
			t.Errorf("%s is written by provisioning but .gitignore admits it", path)
		}
	}
}

func TestGitignoreLeavesTheBootstrapWrappersAddable(t *testing.T) {
	var asked []string
	for _, wrapper := range bootstrapWrappers {
		asked = append(asked, ".githooks/"+wrapper.name)
	}
	for path, ignored := range ignoredIn(t, asked) {
		if ignored {
			t.Errorf("%s is tracked but .gitignore excludes it", path)
		}
	}
}
