# Release Process

This guide describes how to publish a new version of `uit`.

---

## 1. Prepare a Release Pull Request

Create a PR titled like:

```
vYYYY.MM.DD release
```

Include:

- Updates to `README.md` if needed
- A summary of changes (changelog) in the PR description or commit
- Any final cleanups before release

 Once the PR is approved and merged to `main`, proceed to the next step.

You can optionally create a placeholder release commit (without changes) using:

```sh
./scripts/empty-release-commit.sh vYYYY.MM.DD
```

This is useful when you want to reserve a version number or clearly mark the release point.

---

## 2. Tag the Release

After merging to `main`, create a version tag:

```sh
git checkout main
git pull origin main
git tag v2025.03.28
git push origin v2025.03.28
```

This triggers GitHub Actions to:

- Build binaries for multiple platforms
- Package them into `.tar.gz` archives (includes `README.md`)
- Upload the archives to a new GitHub Release

Alternatively, the same steps can be done with:

```sh
./scripts/release.sh vYYYY.MM.DD
```

This script will:

- Ensure you’re on the `main` branch with a clean working directory
- Show the changelog since the last tag
- Create and push the tag
- Trigger the GitHub release workflow

---

## 3. Verify

After the GitHub Actions workflow completes, confirm:

- The release appears under [Releases](https://github.com/mnishiguchi/uit/releases)
- Each `.tar.gz` archive includes:
  - The `uit` binary
  - `README.md`

---

## 4. Unrelease (if needed)

To delete a mistakenly pushed tag:

```sh
git tag -d v2025.03.28-X
git push origin :refs/tags/v2025.03.28-X
```

Then delete the corresponding release on GitHub manually.

---

Happy releasing!
