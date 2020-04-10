from jinja2 import Template
import subprocess
import shlex

subprocess.run(["git", "pull"])

version_tag = subprocess.check_output(shlex.split("git log -1 --pretty=%h"), encoding='utf-8')[0:-1]
project_id = subprocess.check_output(shlex.split("gcloud config get-value project"), encoding='utf-8')[0:-1]

instance_name = subprocess.check_output(shlex.split('gcloud sql instances list --format="value(name)"'), encoding='utf-8')[0:-1]

read_instance_cmd = 'gcloud sql instances describe {} --format="value(connectionName)"'.format(str(instance_name))
connection = subprocess.check_output(shlex.split(read_instance_cmd), encoding='utf-8')[0:-1]

deployment_template = Template(open('./scripts/deployment_template.yaml').read())

deployment = deployment_template.render(VERSION_TAG=version_tag, PROJECT_ID=project_id, DB_CONNECTION=connection)
with open("./deployment.yaml", 'w') as file:
    file.write(deployment)

deployment_script = "#!/bin/bash\n" \
                    "echo Deploying version with tag: {1}\n"\
                    "docker build -t gcr.io/{0}/jorel:{1} .\n" \
                    "docker push gcr.io/{0}/jorel:{1}\n" \
                    "kubectl apply -f deployment.yaml".format(project_id, version_tag)

with open("./run_deployment.sh", 'w') as file:
    file.write(deployment_script)

subprocess.run(["chmod", "u+x", "./run_deployment.sh"])