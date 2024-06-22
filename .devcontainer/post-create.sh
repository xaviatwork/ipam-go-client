#!/usr/bin/env bash
sudo apt-get install make --yes
grep --quiet --fixed-strings --line-regexp 'source .devcontainer/git-completion.bash' ~/.bashrc || echo 'source .devcontainer/git-completion.bash' >> ~/.bashrc
