# OctoActCreateRelease

Create a release with the contents of a given directory.

Contract: you give me the release details, auth token, and a folder with the release assets (only via CLI flags) and I take care of the rest.

Motivation: I wanted a minimal heper to create a release from GitHub Actions.

## Example

```sh
mkdir assets
echo DATA >> assets/data.txt

octoact_create_release \
  -auth-token $GITHUB_TOKEN \
  -repo-owner acorello \
  -repo OctoActCreateRelease \
  -release-name TEST1 \
  -tag-name REL-TEST1 \
  -assets-dir assets
```
