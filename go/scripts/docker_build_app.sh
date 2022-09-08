#!/bin/bash

BASE_VER=0.2
APP_VER=0.4

cp -f ../Dockerfile ../Dockerfile-bkup
cp -f ../Dockerfile-app ../Dockerfile

docker build \
    --build-arg ssh_pub_key="ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC+MUkTMHfVBLcA0TgRMDmhKH5kjuDX2nX8hVMcuu4yJMNsc+eG7j+n2fCs87pobg+E2QxIrVLeKF5fwmdKpfThNV+xPLSdSQdSWSQbXDvPeMzotO77LAxOwquAR3vePrSfj7sNXsP+hg1fSluYaC32PlFZB7WvoiZ6VNF0hnaapdlt4NQZI5ZZGOnKuU1ytXoa0OC/sIyHgEvnB4Yk3k/PzUJuFoOjvM9B/CLsI60vrrARQGQn9WJSbpKi87IfdtH54iWAvGx24DiRukLl60MyeojyByMpNRX93mHPtxjLmfsep7vjoBTti/2y6oqVkpzpPAZc4T2GEyPnWFnli/5r3ZNNTpKtT/XlYaZLOhiIwC5OzuzlpUrbvaCeWlIgxktrc4ijb3NLFlzXYFp01RJ/rmoqHKsmbT55pGWQ6pCcZmZoncsoHvOgX3GSGMRJvMJjMUIDr8in1N1IqgqvEBokhxk+u055djuz17Jlojq44UnylbXSRQziymXTNJTmYEk= jenkins@a7131a1767c9" \
    --build-arg BUILD_IMAGE=lyw1217/ojeommu-base:${BASE_VER} \
    -t lyw1217/ojeommu:${APP_VER} "../"