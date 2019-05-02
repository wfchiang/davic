import os

GOPATH = os.environ['GOPATH']
GO_SRC_DIR = GOPATH + '/src'

# copy the go source 
CMD_COPY_SOURCE = 'cp -rf ./wfchiang ' + GO_SRC_DIR
print (CMD_COPY_SOURCE)
os.system(CMD_COPY_SOURCE)

# build go package 
os.chdir(GOPATH)
CMD_GO_BUILD = 'go build wfchiang/davic'
print (CMD_GO_BUILD)
os.system(CMD_GO_BUILD)
