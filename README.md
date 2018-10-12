# casher

Heavily based on [Travis's CA$H€R](https://github.com/travis-ci/casher)

## Purpose

Casher is a Go tool used to fetch, create, and update caches in an ephemeral CI/CD environment (Jenkins on K8s) using S3.

Casher checks the URLs provided to it for existing caches at the location and stops when it finds one. It will `fetch` the cache and expand it at the root of the OS.

If a cache is not found at any of the URLs, Casher creates one. Casher will `add` the directories specified to be cached and then `push` the cache to a URL.

Casher checks for changes in existing caches using by recursively computing the MD5 checksum for every file in a directory. If a difference in the MD5 checksums are found, casher will pack a new archive and `push` it to a URL.

## Usage

At the start of a CI/CD pipeline:

```
casher fetch -b my-cicd-bucket -k cache/my-app/PR-1 -k cache/my-app/master
```

Adding directories to cache:

```
casher add -p /home/jenkins/.ivy2 -p /home/jenkins/.sbt/boot
```

Saving cache after build pipeline:

```
casher push -b my-cicd-bucket -k cache/my-app/PR-1
```

### Jenkins Pipeline Example

Having a simple `withCaching.groovy` script that wraps a body of work by fetching & adding before and pushing any changes after.

```groovy
#!/usr/bin/groovy

def call(Map config = [:], body) {
  def awsRegion = config.get('awsRegion', 'us-east-1')
  def cacheBucket = config.get('cacheBucket', 'my-cicd-bucket')
  def cacheDirectories = config.get('cacheDirectories', [])
  def jobBaseName = "${env.JOB_NAME}".split('/').tail().join("/")

  if (cacheDirectories.size > 0) {
    def cacheDirs = cacheDirectories.join(' -p ')

    stage('loading cache') {
      container('cacher') {
        withAWSRole() {
          sh "casher fetch -r ${awsRegion} -b ${cacheBucket} -k cache/${jobBaseName}"
          sh "casher add -p ${cacheDirs}"
        }
      }
    }

    body()

    stage('store build cache') {
      container('cacher') {
        withAWSRole() {
          sh "casher push -r ${awsRegion} -b ${cacheBucket} -k cache/${jobBaseName}"
        }
      }
    }
  } else {
    body()
  }
}
```

Example usage:

```groovy
container('scala') {
  withCaching(cacheDirectories: ['/root/.sbt/boot', '/root/.ivy2']) {
    sh "sbt clean test"
  }
}
```
