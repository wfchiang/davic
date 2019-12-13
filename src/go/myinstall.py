import os
import sys

CMD_COPY   = 'cp -rf'
CMD_REMOVE = 'rm -rf'
GOPATH = os.environ['GOPATH']
GO_SRC_DIR = os.path.abspath(GOPATH + '/src')

PACKAGE_NAME = 'wfchiang/davic'

# check and setup the properties 
print ('==== Check and Setup the Properties ====')
CMD_RM_SOURCE = ''
CMD_COPY_SOURCE = ''

if (sys.platform in ['win32', 'win64']): 
    CMD_COPY = 'xcopy /E /I'
    CMD_REMOVE = 'del'
    CMD_RM_SOURCE = CMD_REMOVE + ' ' + os.path.abspath(GO_SRC_DIR+'/'+PACKAGE_NAME)
    CMD_COPY_SOURCE = CMD_COPY + ' ' + os.path.abspath('./wfchiang/*') + ' ' + os.path.abspath(GO_SRC_DIR+'/wfchiang')
else: 
    CMD_RM_SOURCE = CMD_REMOVE + ' ' + os.path.abspath(GO_SRC_DIR+'/'+PACKAGE_NAME)
    CMD_COPY_SOURCE = CMD_COPY + ' ' + os.path.abspath('./wfchiang') + ' ' + GO_SRC_DIR

print ('copy command: ' + CMD_COPY)
print ('remove command: ' + CMD_REMOVE)
print ('copy source command: ' + CMD_COPY_SOURCE)
print ('remove source command: ' + CMD_RM_SOURCE)
print ('')

# remove the source folder 
print ('==== Remove the Source Folder ====')
print (CMD_RM_SOURCE)
os.system(CMD_RM_SOURCE)
print ('')

# copy the go source 
print ('==== Copy the Source Files ====')
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