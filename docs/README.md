# Build the documentation

The documentation is built using Github pages on each commit to the main branch.

## Build locally

Tested on Fedora 35.

Install the dependciesm make sure to install ruby < 3.0:
```
dnf module install ruby:2.7/default
dnf install ruby-devel
```

Install jekyll and all dependencies:
```
bundle install --path vendor
```

Build and serve the site locally:
```
bundle exec jekyll serve
```

bundle config set --local path 'vendor/bundle'
