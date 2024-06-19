package osfilesys

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"homework/internal/domain/file"
	"io"
	"os"
	"path"
	"testing"
	"time"
)

type testFile struct {
	data []byte
	name string
}

type repoSuite struct {
	suite.Suite
	repo     file.Repository
	rootDir  string
	fileList []testFile
	modTime  int64
}

func (suite *repoSuite) SetupSuite() {
	curDir, _ := os.Getwd()
	suite.rootDir = path.Join(curDir, "test_dir")
	suite.fileList = []testFile{
		{[]byte("first"), "first.txt"},
		{[]byte("second"), "second.txt"},
		{[]byte("third"), "third.txt"},
	}
	_ = os.Mkdir(suite.rootDir, 0777)
	suite.modTime = time.Now().Unix()
	for _, v := range suite.fileList {
		_ = os.WriteFile(path.Join(suite.rootDir, v.name), v.data, 0777)
	}
	suite.repo, _ = NewFileRepository(suite.rootDir)
}

func (suite *repoSuite) TearDownSuite() {
	_ = os.RemoveAll(suite.rootDir)
}

func (suite *repoSuite) TestNewFileRepository() {
	suite.Run("valid directory path", func() {
		_, err := NewFileRepository(suite.rootDir)
		assert.NoError(suite.T(), err)
	})
	suite.Run("invalid directory path", func() {
		_, err := NewFileRepository(suite.rootDir + "wrong")
		assert.Error(suite.T(), err)
	})
}

func (suite *repoSuite) TestFileRepository_GetList() {
	list, err := suite.repo.GetList()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(suite.fileList), len(list))
	for i := range list {
		assert.Equal(suite.T(), suite.fileList[i].name, list[i].Name)
		assert.Equal(suite.T(), suite.modTime, list[i].ModTime.Unix())
		assert.Equal(suite.T(), len(suite.fileList[i].data), int(list[i].Size))
	}
}

func (suite *repoSuite) TestFileRepository_GetInfo() {
	suite.Run("valid files", func() {
		for _, f := range suite.fileList {
			info, err := suite.repo.GetInfo(f.name)
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), f.name, info.Name)
			assert.Equal(suite.T(), len(f.data), int(info.Size))
			assert.Equal(suite.T(), suite.modTime, info.ModTime.Unix())
		}
	})
	suite.Run("invalid files' names", func() {
		_, err := suite.repo.GetInfo("baraboba")
		assert.Error(suite.T(), err)
	})
}

func (suite *repoSuite) TestFileRepository_GetData() {
	for _, f := range suite.fileList {
		data, err := suite.repo.GetData(f.name)
		assert.NoError(suite.T(), err)
		b, _ := io.ReadAll(data)
		assert.Equal(suite.T(), string(f.data), string(b))
		_ = data.Close()
	}
}

func TestFileRepository(t *testing.T) {
	suite.Run(t, new(repoSuite))
}
