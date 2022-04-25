import os
import yaml
import sys

os.makedirs('plugins/ImageStream', exist_ok=True)

if os.path.exists("index.yaml"):
    file = open("index.yaml","r")
    not_present = 1
    index = yaml.safe_load(file)
    for plugin in index['plugins']:
        if plugin['name'] == 'ImageStreamPlugin':
            not_present = 0
            break
    if not_present:
        index['plugins'].append({"name": "ImageStreamPlugin", "path": "https://github.com/%s/crane-plugins/raw/main/plugins/ImageStream/index.yaml"%sys.argv[2]})
    file = open("index.yaml","w")
    yaml.dump(index, file)
    file.close()

else:
    file = open("index.yaml","a+")

    index = yaml.safe_load(file)

    index = {}
    index['kind'] = 'PluginIndex'
    index['apiVersion'] = 'crane.konveyor.io/v1alpha1'
    index['plugins'] = []

    index['plugins'].append({"name": "ImageStreamPlugin", "path": "https://github.com/%s/crane-plugins/raw/main/plugins/ImageStream/index.yaml"%sys.argv[2]})

    yaml.dump(index, file)
    file.close()

# create or append in plugin index
if os.path.exists('plugins/ImageStream/index.yaml'):

    file = open("plugins/ImageStream/index.yaml","r")

    index = yaml.safe_load(file)

    index['versions'].append({})
    index['versions'][-1] = {
        'name': 'ImageStreamPlugin',
        'shortDescription': 'ImageStreamPlugin',
        'description': 'this is ImageStreamPlugin',
        'version': sys.argv[1],
        'binaries': [
            {
                'os': 'linux',
                'arch': 'amd64',
                'uri': "https://github.com/%s/releases/download/%s/amd64-linux-imagestreamplugin-%s"%(sys.argv[3], sys.argv[1],sys.argv[1]),
            },
            {
                'os': 'darwin',
                'arch': 'amd64',
                'uri': "https://github.com/%s/releases/download/%s/amd64-darwin-imagestreamplugin-%s"%(sys.argv[3], sys.argv[1],sys.argv[1]),
            },
            {
                'os': 'darwin',
                'arch': 'arm64',
                'uri': "https://github.com/%s/releases/download/%s/arm64-darwin-imagestreamplugin-%s"%(sys.argv[3], sys.argv[1],sys.argv[1]),
            },
        ],
        'optionalFields': [
            { 
                'flagName': "registry-replacement",
                'help':     "Map of image registry paths to swap on transform, in the format original-registry1=target-registry1,original-registry2=target-registry2...",
                'example':  "docker-registry.default.svc:5000=image-registry.openshift-image-registry.svc:5000,docker.io/foo=quay.io/bar",
            },
        ]
    }

    file = open("plugins/ImageStream/index.yaml","w")

    yaml.dump(index, file)
    file.close()
    
else:
    file = open("plugins/ImageStream/index.yaml","a+")

    index = yaml.safe_load(file)

    index = {}
    index['kind'] = 'Plugin'
    index['apiVersion'] = 'crane.konveyor.io/v1alpha1'
    index['versions'] = []

    index['versions'].append({})
    index['versions'][0] = {
        'name': 'ImageStreamPlugin',
        'shortDescription': 'ImageStreamPlugin',
        'description': 'this is ImageStreamPlugin',
        'version': sys.argv[1],
        'binaries': [
            {
                'os': 'linux',
                'arch': 'amd64',
                'uri': "https://github.com/%s/releases/download/%s/amd64-linux-imagestreamplugin-%s"%(sys.argv[3], sys.argv[1],sys.argv[1]),
            },
            {
                'os': 'darwin',
                'arch': 'amd64',
                'uri': "https://github.com/%s/releases/download/%s/amd64-darwin-imagestreamplugin-%s"%(sys.argv[3], sys.argv[1],sys.argv[1]),
            },
            {
                'os': 'darwin',
                'arch': 'arm64',
                'uri': "https://github.com/%s/releases/download/%s/arm64-darwin-imagestreamplugin-%s"%(sys.argv[3], sys.argv[1],sys.argv[1]),
            },
        ],
        'optionalFields': [
            {  
                'flagName': "registry-replacement",
                'help':     "Map of image registry paths to swap on transform, in the format original-registry1=target-registry1,original-registry2=target-registry2...",
                'example':  "docker-registry.default.svc:5000=image-registry.openshift-image-registry.svc:5000,docker.io/foo=quay.io/bar",
            },
        ]
    }
    
    yaml.dump(index, file)
    file.close()