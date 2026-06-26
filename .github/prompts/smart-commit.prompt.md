---
description: "Stage unstaged files, review changes, generate a meaningful commit message, and commit"
mode: "agent"
tools: ['terminal', 'codebase', 'execute']
---

# Smart Commit Workflow

Please perform the following git workflow:

## Step 1: Check Status
Run `git status` to identify all unstaged files and changes.

## Step 2: Review Changes
Use `git diff` to review all unstaged changes. Read through the changes carefully to understand:
- What files were modified
- What functionality was added, changed, or removed
- The context and purpose of the changes

## Step 3: Stage Changes
Stage all unstaged files using `git add .` (or selectively stage specific files if needed).

## Step 4: Generate Commit Message
Based on your review of the changes, generate a meaningful commit message that follows these best practices:

- **Format**: Use conventional commits format when appropriate (feat:, fix:, docs:, refactor:, test:, chore:, etc.)
- **Subject line**: Clear, concise summary (50 chars or less) in imperative mood
- **Body** (if needed): Explain what and why, not how. Wrap at 72 characters.
- **Be specific**: Reference specific components, features, or fixes
- **Group related changes**: If multiple related changes, describe the overall goal

### Commit Message Template:
```
<type>(<scope>): <subject>

<body>

<footer>
```

## Step 5: Commit
Commit the staged changes with the generated commit message using `git commit -m "message"` (or `git commit` for multi-line messages).

## Step 6: Confirm
Show the commit summary with `git log -1` to confirm the commit was successful.

---

**Note**: If there are no unstaged changes, inform the user. If you need clarification about the purpose of certain changes, ask before generating the commit message.
