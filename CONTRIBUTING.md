# Contributing to Janus

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

The following is a set of guidelines for contributing to Janus and its packages, 
which are hosted in the [HelloFresh Organization](https://github.com/hellofresh) on GitHub.
These are just guidelines, not rules. Use your best judgment, and feel free to propose changes 
to this document in a pull request.

## Code of Conduct

This project adheres to the Contributor Covenant [code of conduct](CODE_OF_CONDUCT.md).
By participating, you are expected to uphold this code.
Please report unacceptable behavior to [engineering@hellofresh.com](mailto:engineering@hellofresh.com).

We accept contributions via Pull Requests on [Github](https://github.com/hellofresh/janus).

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report for Janus. Following these guidelines helps maintainers 
and the community understand your report :pencil:, reproduce the behavior :computer: :computer:, and find related 
reports :mag_right:.

Before creating bug reports, please check if the bug was already reported before as you might find out that you don't 
need to create one. When you are creating a bug report, please [include as many details as possible](#how-do-i-submit-a-good-bug-report). 

#### How Do I Submit A (Good) Bug Report?

Bugs are tracked as [GitHub issues](https://guides.github.com/features/issues/). Create an issue on provide the following information.

Explain the problem and include additional details to help maintainers reproduce the problem:

* **Use a clear and descriptive title** for the issue to identify the problem.
* **Describe the exact steps which reproduce the problem** in as many details as possible. For example, start by explaining how you started Janus, 
e.g. which command exactly you used in the terminal. When listing steps, **don't just say what you did, but explain how you did it**.
* **Provide specific examples to demonstrate the steps**. Include links to files or GitHub projects, or copy/pasteable snippets, which you use in those examples. 
If you're providing snippets in the issue, use [Markdown code blocks](https://help.github.com/articles/markdown-basics/#multiple-lines).
* **Describe the behavior you observed after following the steps** and point out what exactly is the problem with that behavior.
* **Explain which behavior you expected to see instead and why.**

Include details about your configuration and environment:

* **Which version of Janus are you using?**
* **What's the name and version of the OS you're using**?

### Your First Code Contribution

Unsure where to begin contributing to Janus? You can start by looking through these `beginner` and `help-wanted` issues:

* [Beginner issues][beginner] - issues which should only require a few lines of code, and a test or two.
* [Help wanted issues][help-wanted] - issues which should be a bit more involved than `beginner` issues.

Both issue lists are sorted by total number of comments. While not perfect, number of comments is a reasonable proxy for impact a given change will have.

### Pull Requests

* Include screenshots and animated GIFs in your pull request whenever possible.
* Follow the [Go](https://github.com/golang/go/wiki/CodeReviewComments) styleguides.
* Include thoughtfully-worded, well-structured tests.
* Document new code
* End files with a newline.


Happy Coding from the HelloFresh Engineering team!


## Compiling and Debugging

### Compiling
After cloning the solution, using the `master` branch, the code should be compiling at all times. If it's not, please file a bug.

### Debugging
To quickly start a debug session, a simple api has been provided on the examples folder. To use it, a `janus.toml` file at root of the cloned directory and paste the following content to it and save

```
################################################################
# Global configuration
################################################################
port = 8080

[log]
  level = "debug"

################################################################
# API configuration backend
################################################################
[web]
  port = 8081

  [web.credentials]
    secret = "secret"

    [web.credentials.basic]
    users = {admin = "admin"}

[database]
  dsn = "file://{Path To Cloned Directory}}/examples/front-proxy-EchoApi/"
```

If you're using VSCode, you can use the following settings on the `launch.json`:

```
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch",
            "type": "go",
            "args": ["start"],
            "request": "launch",
            "program": "${workspaceFolder}"
        }
    ]
}
```

When you hit debug, the application will start, bind to port 8080. When access `http://localhost:8080/echo` you should get a echo back from the publically available postman-echo service (https://learning.postman.com/docs/developer/echo-api/).