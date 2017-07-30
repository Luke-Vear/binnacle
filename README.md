<p align="center">
    <img src="https://upload.wikimedia.org/wikipedia/commons/thumb/c/c4/Compass_rose_Cantino.svg/200px-Compass_rose_Cantino.svg.png" height="200" />
</p>

# binnacle

As nodes in your Kubernetes cluster are terminated, scaled or suffer failure, your ingress controllers will be rescheduled and their external IP addresses will change.

Binnacle is an in cluster DNS manager that keeps your DNS provider pointing at your ingress controllers. Once configured, binnacle polls the Kubernetes API for changes to your ingress, then if there are any changes, updates your records in DNS.

## Setup

Instructions to follow.

## DNS Providers

Binnacle currently works with AWS route53.

# Kubernetes are working on a native solution, so you should probably use theirs:

https://github.com/kubernetes-incubator/external-dns
