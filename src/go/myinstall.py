import os

GOPATH = os.environ['GOPATH']
GO_SRC_DIR = GOPATH + '/src'

PACKAGE_NAME = "wfchiang/davic"

# remove the source folder 
print ('==== Remove the Source Folder ====')
CMD_RM_SOURCE = 'rm -rf ' + GO_SRC_DIR + '/' + PACKAGE_NAME
print (CMD_RM_SOURCE)
os.system(CMD_RM_SOURCE)
print ('')

# copy the go source 
print ('==== Copy the Source Files ====')
CMD_COPY_SOURCE = 'cp -rf ./wfchiang ' + GO_SRC_DIR
print (CMD_COPY_SOURCE)
os.system(CMD_COPY_SOURCE)
print ('')

# change dir
print ('==== Change the Current Dir ====')
print ('os.chdir(' + GOPATH + ')')
os.chdir(GOPATH)
print ('')

# uninstall the package 
print ('==== Uninstall the Package ====')
CMD_GO_UNINSTALL = 'go clean -i ' + PACKAGE_NAME
print (CMD_GO_UNINSTALL)
os.system(CMD_GO_UNINSTALL)
print ('')

# build go package 
print ('==== Install the Package ====')
CMD_GO_BUILD = 'go build ' + PACKAGE_NAME
print (CMD_GO_BUILD)
os.system(CMD_GO_BUILD)
print ('')

# unit testing 
print ('==== Unit Tests ====')
CMD_UNIT_TEST = 'go test ' + PACKAGE_NAME
print (CMD_UNIT_TEST)
os.system(CMD_UNIT_TEST)
print ('')