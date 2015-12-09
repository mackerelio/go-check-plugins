#!/bin/sh
set -e

deploykey=~/.ssh/deploy.key

echo "
Host github.com
    StrictHostKeyChecking no
    IdentityFile $deploykey
" >> ~/.ssh/config
openssl aes-256-cbc -K $encrypted_87d0e2b1ee75_key -iv $encrypted_87d0e2b1ee75_iv -in tool/travis/go-check-plugins.pem.enc -out $deploykey -d
chmod 600 $deploykey
git config --global user.email "mackerel-developers@hatena.ne.jp"
git config --global user.name  "mackerel"
git remote set-url origin git@github.com:mackerelio/go-check-plugins.git
tool/autotag
