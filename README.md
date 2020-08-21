![GitHub](https://img.shields.io/github/license/mjpitz/credentials-operator.svg)
<!--![branch](https://github.com/mjpitz/credentials-operator/workflows/branch/badge.svg?branch=main)-->
[![Go Report Card](https://goreportcard.com/badge/github.com/mjpitz/credentials-operator)](https://goreportcard.com/report/github.com/mjpitz/credentials-operator)
![Google Analytics](https://www.google-analytics.com/collect?v=1&cid=555&t=pageview&ec=repo&ea=open&dp=credentials-operator&dt=credentials-operator&tid=UA-172921913-1)

# credentials-operator

Easily generate, store, and share credentials with applications running within a Kubernetes cluster.

In working on Kubernetes, I often find the process of generating credentials to be tedious.
Beyond generation, these credentials often need to be shared between systems.
For example, an application using MySQL needs to know the username and password it should use to connect.

## Status

Mostly a toy.
Wanted to play around with Ranchers' [Wrangler](https://github.com/rancher/wrangler) system and this seemed like a nifty idea.
If interest grows, I'm willing to support this a bit more.

## Getting Started

### Installing the Operator

```bash
git clone git@github.com:mjpitz/credentials-operator.git
helm upgrade -i credentials-operator ./charts/credentials-operator
```

### Creating a Credential

The configuration below can be used in conjunction with [Bitnami's MySQL chart](https://github.com/bitnami/charts/tree/master/bitnami/mysql).

```bash
cat << 'EOF' | kubectl apply -f -
apiVersion: credentials.mjpitz.com/v1alpha1
kind: Credential
metadata:
  name: bitnami-mysql-passwords
spec:
  credentials:
    - key: DB_ROOT_PASSWORD
      requirements:
        length: 10
        characterSet: a-zA-Z0-9
    - key: DB_PASSWORD
      requirements:
        length: 10
        characterSet: a-zA-Z0-9
    - key: DB_REPLICATION_PASSWORD
      requirements:
        length: 10
        characterSet: a-zA-Z0-9
  views:
    - secretRef:
        name: bitnami-mysql-secret
      stringDataTemplate:
        mysql-root-password: ${DB_ROOT_PASSWORD}
        mysql-password: ${DB_PASSWORD}
        mysql-replication-password: ${DB_REPLICATION_PASSWORD}
    - secretRef:
        name: myapp-db-config
      stringDataTemplate:
        DB_DRIVER: mysql
        DB_CONNECTION_STRING: myapp:${DB_PASSWORD}@tcp(mysql:3306)/dbname
EOF
```

This credential creates a total of **3** secrets.
The first, `bitnami-mysql-passwords` is used to store the generated credentials in the `.spec.credentials` block. 
The second and third, (`bitnami-mysql-secret` and `myapp-db-config` respectively) are synthesized by replacing environment variables in the `stringDataTemplate` with associated values.
For example, the YAML below shows how this configuration renders into associated manifests. 

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: bitnami-mysql-secret
stringData:
  DB_ROOT_PASSWORD: adminpass
  DB_PASSWORD: notsecure
  DB_REPLICATION_PASSWORD: replicationpass
---
apiVersion: v1
kind: Secret
metadata:
  name: bitnami-mysql-secret
stringData:
  mysql-root-password: adminpass
  mysql-password: notsecure
  mysql-replication-password: replicationpass
---
apiVersion: v1
kind: Secret
metadata:
  name: myapp-db-config
stringData:
  DB_DRIVER: mysql
  DB_CONNECTION_STRING: myapp:notsecure@tcp(mysql:3306)/dbname
```
