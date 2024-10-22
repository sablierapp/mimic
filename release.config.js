module.exports = {
  "branches":  [
    { "name": "main" },
  ],
  "plugins": [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/github",
    ["@semantic-release/exec", {
      "publishCmd": "echo \"${nextRelease.notes}\" > /tmp/release-notes.md; goreleaser release --release-notes /tmp/release-notes.md --clean"
    }]
  ]
}