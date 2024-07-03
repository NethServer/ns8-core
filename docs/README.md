# Build the documentation

The documentation is built using Github pages on each commit to the main branch.

## Build locally

Tested on Fedora 36 and Ubuntu 22.04.
Install the dependencies make sure to install ruby >= 3.0.

- Fedora:
    ```
    dnf module install ruby
    dnf install ruby-devel
    ```

- Ubuntu:
    ```
    apt-get install ruby-full
    ```

Install jekyll and all dependencies:
```
bundle config set --local path 'vendor/bundle'
bundle install --path vendor
```

Build and serve the site locally:
```
bundle exec jekyll serve
```

Add `-B` to run in detached mode.

To kill the process:
```
kill -9 [processID]
```
