package lib

import (
	"fmt"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var specChineseCreateSymlink = SpecText{

	synopsisText: "创建符号链接",

	paramText: "[sym_url] [target_url] [options]",

	syntaxText: ` 
    ossutil create-symlink sym_url target_object [--encoding-type url] 
`,

	detailHelpText: ` 
    该命令在oss上创建符号链接文件，链接的目标文件必须为相同bucket下的文件，且文件类型非符
    号链接。即，symlink_url必须为形如oss://bucket/object的cloud_url，target_object为object名。

    创建符号链接时：
        不检查目标文件是否存在，
        不检查目标文件类型是否合法，
        不检查目标文件是否有权限访问， 
        以上检查，都推迟到GetObject等需要访问目标文件的API。
    如果试图添加的文件已经存在，并且有访问权限。新添加的文件将覆盖原来的文件。

    通过stat命令可以查看符号链接的目标文件。

    更多信息见官网API文档：https://help.aliyun.com/document_detail/45126.html?spm=5176.doc31979.6.870.x3Tqsh


用法：

    ossutil create-symlink oss://bucket/symlink-object target-object
`,

	sampleText: ` 
    ossutil create-symlink oss://bucket1/object1 object2 
        创建从指向object2的符号链接object1。
`,
}

var specEnglishCreateSymlink = SpecText{

	synopsisText: "Create symlink of object",

	paramText: "[sym_url] [target_url] [options]",

	syntaxText: ` 
    ossutil create-symlink sym_url target_object [--encoding-type url]
`,

	detailHelpText: ` 
    The command create symlink of object in oss, the target object must be object in the same 
    bucket of symlink object, and the file type of target object must not be symlink. So, 
    symlink_url must be in format: oss://bucket/object, and target_object is the object name 
    of target object.  

    When create symlink:
        Will not check whether target object exists;
        Will not check whether target object type is valid;
        Will not check whether if have access permission of target object.
    The check will be done when visiting GetObject, etc.

    If the symlink object exist, and has access permission, the object newly created will cover 
    old object.

    We can use stat command to query the target object of symlink object.

    More information about symlink see: https://help.aliyun.com/document_detail/45126.html?spm=5176.doc31979.6.870.x3Tqsh


Usage:

    ossutil create-symlink oss://bucket/symlink-object target-object
`,

	sampleText: ` 
    ossutil create-symlink oss://bucket1/object1 object2 
        Create symlink object named object1, which point to object2.
`,
}

// CreateSymlinkCommand is the command list buckets or objects
type CreateSymlinkCommand struct {
	command Command
}

var createSymlinkCommand = CreateSymlinkCommand{
	command: Command{
		name:        "create-symlink",
		nameAlias:   []string{},
		minArgc:     2,
		maxArgc:     2,
		specChinese: specChineseCreateSymlink,
		specEnglish: specEnglishCreateSymlink,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
			OptionEncodingType,
			OptionConfigFile,
			OptionEndpoint,
			OptionAccessKeyID,
			OptionAccessKeySecret,
			OptionSTSToken,
			OptionRetryTimes,
		},
	},
}

// function for FormatHelper interface
func (cc *CreateSymlinkCommand) formatHelpForWhole() string {
	return cc.command.formatHelpForWhole()
}

func (cc *CreateSymlinkCommand) formatIndependHelp() string {
	return cc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (cc *CreateSymlinkCommand) Init(args []string, options OptionMapType) error {
	return cc.command.Init(args, options, cc)
}

// RunCommand simulate inheritance, and polymorphism
func (cc *CreateSymlinkCommand) RunCommand() error {
	encodingType, _ := GetString(OptionEncodingType, cc.command.options)
	cloudURL, err := CloudURLFromString(cc.command.args[0], encodingType)
	if err != nil {
		return err
	}

	targetURL, err := StorageURLFromString(cc.command.args[1], encodingType)
	if err != nil {
		return err
	}

	targetObject := targetURL.ToString()
	if targetURL.IsCloudURL() {
		if targetURL.(CloudURL).bucket != cloudURL.bucket {
			return fmt.Errorf("the bucket of target object: %s must be the same with the bucket of symlink object: %s", targetURL.(CloudURL).bucket, cloudURL.bucket)
		}
		targetObject = targetURL.(CloudURL).object
	}

	bucket, err := cc.command.ossBucket(cloudURL.bucket)
	if err != nil {
		return err
	}

	return cc.ossCreateSymlink(bucket, cloudURL.object, targetObject)
}

func (cc *CreateSymlinkCommand) ossCreateSymlink(bucket *oss.Bucket, symlinkObject, targetObject string) error {
	retryTimes, _ := GetInt(OptionRetryTimes, cc.command.options)
	for i := 1; ; i++ {
		err := bucket.PutSymlink(symlinkObject, targetObject)
		if err == nil {
			return err
		}
		if int64(i) >= retryTimes {
			return ObjectError{err, bucket.BucketName, symlinkObject}
		}
	}
}