pipeline {
  agent any

  stages {
    stage('PR - Feature - Teste Unitario') {
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
      deleteDir()
    }

    success {
      echo "Release Success"
    }

    failure {
      echo "Release Failed"
    }
  }
}
