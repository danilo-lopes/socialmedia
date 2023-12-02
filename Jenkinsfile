pipeline {
  agent any

  stages {
    stage('PR - Teste Unitario') {
      when { branch 'PR-*' }

      steps {
        script {
          sh '''
            echo Este step é executado apenas quando tiver um PR aberto
          '''
        }
      }
    }
  }

  post {
    always {
      node('main'){
        deleteDir()
      }
    }

    success {
      echo "Release Success"
    }

    failure {
      echo "Release Failed"
    }
  }
}
