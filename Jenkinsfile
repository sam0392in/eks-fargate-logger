import org.jenkinsci.plugins.pipeline.modeldefinition.Utils

// Your config here:
def AWS_PROFILE_DEV = '123445566778'
def AWS_PROFILE_LIVE = '998877665544'
def AWS_ECR_NAME = 'eks-fargate-logger'

def label = "devops-${UUID.randomUUID().toString()}"

def normalizeBranchName(String branchName) {
        return branchName.replaceAll(/[^a-zA-Z0-9]/, '-').replaceAll(/[\-]+/, '-').toLowerCase()
    }

def compileCode () {
	sh("go build -x -v main.go")
}

def buildImage() {
  sh """
    echo "Preparing image for aggregators"
    /kaniko/executor -f ./Dockerfile \
    -c . \
    --cache=true \
    --cache-repo=${env.AWS_PROFILE_DEPLOY}.dkr.ecr.eu-west-1.amazonaws.com/${env.AWS_ECR_NAME} \
    --registry-mirror docker-registry.sam.io \
    --registry-mirror registry-1.docker.io \
    --destination=${env.AWS_PROFILE_DEPLOY}.dkr.ecr.eu-west-1.amazonaws.com/${env.AWS_ECR_NAME}:${env.TAG}

    rm -rf /kaniko/0
    mkdir -p /workspace
  """
}

def ecrScan(String awsProfile) {
	sh """
		export AWS_PROFILE=${awsProfile} && sam-ecr-scanner scan --repo ${env.AWS_ECR_NAME} --tag ${env.TAG} --severity critical,high
	"""
}

def updateChartImage() {
    sh """
      set -e
      cd k8s/
      sed -i 's/  tag:.*/  tag: ${env.TAG}/' environments/${env.VALUE_FILE}
      cat environments/${env.VALUE_FILE}
    """
}

def bumpChartVersion() {
  sh "helm-update-version -d ./k8s -b minor"
}

def getVersionFromChart(String chartLocation) {
    def version = sh(script: "cat ${chartLocation}/Chart.yaml | grep '^version:' | awk '{ print \$2 }'", returnStdout: true).trim()
    return version
}

def release(String chartName) {
    sh """#!/bin/bash
        set -e
        make release TAG=${env.TAG} APPNAME=${chartName} ENV=${env.ENV}
    """
}

def syncGitChanges(){
    withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'github-app-sam', usernameVariable: 'GITHUB_APP', passwordVariable: 'GITHUB_ACCESS_TOKEN']]) {
        def prepareGitScript = """
        set -e
        git config --global --add safe.directory '*'
        git config --global user.email $GITHUB_APP
        git config --global user.name "samsam"
        git remote set-url origin https://$GITHUB_APP:$GITHUB_ACCESS_TOKEN@github.com/sam/eks-fargate-logger.git
        git fetch --all
        git checkout ${env.BRANCH}
        """
        sh(prepareGitScript)

        def tag = getVersionFromChart("./k8s")

        sh """
          git add k8s && git commit -m "Bumping up helm chart version to: ${tag}"
          git push origin ${env.BRANCH}
        """
    }
}

def getReleaseTag(){
    withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'github-app-sam', usernameVariable: 'GITHUB_APP', passwordVariable: 'GITHUB_ACCESS_TOKEN']]) {
        def prepareGitScript = """
        set -e
        git config --global --add safe.directory '*'
        git config --global user.email $GITHUB_APP
        git config --global user.name "sam"
        git remote set-url origin https://$GITHUB_APP:$GITHUB_ACCESS_TOKEN@github.com/sam/eks-fargate-logger.git
        git fetch --all --tags
        """
        sh(prepareGitScript)
        def currentTag = sh(script: "git for-each-ref --sort=creatordate --format '%(refname)' refs/tags | tail -1 | awk -F '/' '{print \$3}'", returnStdout: true).trim()

        return currentTag
    }

}

def tagRelease(){
    withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'github-app-sam', usernameVariable: 'GITHUB_APP', passwordVariable: 'GITHUB_ACCESS_TOKEN']]) {
        def prepareGitScript = """
        set -e
        git config --global --add safe.directory '*'
        git config --global user.email $GITHUB_APP
        git config --global user.name "samsam"
        git remote set-url origin https://$GITHUB_APP:$GITHUB_ACCESS_TOKEN@github.com/sam/eks-fargate-logger.git
        git fetch --all --tags
        """
        sh(prepareGitScript)
        def currentTag = getReleaseTag()
        def tag = sh(script: "helm-update-version -d ./k8s -v ${currentTag} -c ${env.CHANGE_TYPE}", returnStdout: true).trim()

        if (sh(script: "git tag -l \"$tag\"", returnStdout: true)) {
            println("Tag $tag already exists")
            return
        }
        sh("git tag -a ${tag} -m \"Release version: ${tag}\"")
        sh("git push origin $tag")
    }
}

podTemplate(
  serviceAccount: 'jenkins',label: label,  yaml: readTrusted('build-agent.yaml')) {
  node(label) {
    properties([
      disableConcurrentBuilds()
    ])

    def myRepo = checkout scm
    def gitCommit = myRepo.GIT_COMMIT
    def branchName = env.BRANCH_NAME
    branchName = branchName.replaceAll(/[^a-zA-Z0-9]/, '-').replaceAll(/[\-]+/, '-').toLowerCase()
    env.BRANCH = branchName
    def shortGitCommit = sh returnStdout: true, script: 'git rev-parse --short HEAD'

    // these env vars are passed to the container
    env.COMMIT_TAG = shortGitCommit.trim()
    env.AWS_ECR_NAME = AWS_ECR_NAME
    env.TAG = "${env.BRANCH}-${env.COMMIT_TAG}"

    def chartName = ""
    def serviceHost = ""
    switch (branchName) {
      case ~/main/:
        chartName = "eks-fargate-logger"
        env.AWS_PROFILE_DEPLOY = AWS_PROFILE_LIVE
        env.VALUE_FILE = "values-live.yaml"
        env.ENV = "live"
        env.CHANGE_TYPE = "minor"
        break;
      case ~/hotfix/:
        chartName = "eks-fargate-logger-hotfix"
        env.AWS_PROFILE_DEPLOY = AWS_PROFILE_LIVE
        env.VALUE_FILE = "values-live.yaml"
        env.ENV = "live"
        env.CHANGE_TYPE = "hotfix"
        break;
      default:
        chartName = "eks-fargate-logger-dev"
        env.AWS_PROFILE_DEPLOY = AWS_PROFILE_DEV
        env.VALUE_FILE = "values-dev.yaml"
        env.ENV = "dev"
        break;
    }

    stage('Compile Code') {
      container('go') {
         compileCode()
      }
    }


    stage('Build Image') {
      container('imagebuilder') {
        if(branchName == "main") {
          sh "export AWS_PROFILE=live"
          buildImage()
          updateChartImage()
        } else if(branchName == "dev") {
          sh "export AWS_PROFILE=default"
          buildImage()
          updateChartImage()
        } else {
          echo "Image build not required for feature branches"
        }
      }
    }

    switch (branchName) {
      case ~/main/:
      case ~/hotfix/:
        stage('Deploy'){
          container('cicd') {
            withCredentials([file(credentialsId: 'devops-live-kube-config', variable: 'KUBECONFIG')]) {
              withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'github-app-sam', usernameVariable: 'GITHUB_APP', passwordVariable: 'GITHUB_ACCESS_TOKEN']]) {
                echo 'Deploying in Live cluster-services'
                def currentTag = getReleaseTag()
                sh(script: "helm-update-version -d ./k8s -v ${currentTag} -c ${env.CHANGE_TYPE}")
                release(chartName)
              }
            }
          }
        }

				stage('ECR Scan'){
				  container('cicd') {
					  ecrScan('live')
				  }
				}
        break

      case ~/dev/:
        stage('Deploy'){
          container('cicd'){
            withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'github-app-sam', usernameVariable: 'GITHUB_APP', passwordVariable: 'GITHUB_ACCESS_TOKEN']]) {
              echo 'Deploying in Dev cluster-services'
              release(chartName)
            }
          }
        }

        stage('ECR Scan'){
          container('cicd') {
            ecrScan('dev')
          }
        }
        break

      default:
        stage('Deploy'){
          container('cicd'){
              echo "Deployment not required for PR and branches other than main and dev"
          }
        }
        break
    }

    stage('Tag Release') {
      container('cicd') {
        if(branchName == "main") {
          tagRelease()
        }
      }
    }

  }
}