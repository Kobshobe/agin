package agin

type Mysql struct {
	Path            string
	Config          string `mapstructure:"config" json:"config" yaml:"config"`
	Dbname          string `mapstructure:"db-name" json:"dbname" yaml:"db-name"`
	Username        string `mapstructure:"username" json:"username" yaml:"username"`
	Password        string `mapstructure:"password" json:"password" yaml:"password"`
	MaxIdleConns    int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"`
	MaxOpenConns    int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"`
	LogMode         bool   `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"`
	LogZap          string `mapstructure:"log-zap" json:"logZap" yaml:"log-zap"`
	LocalPath       string `mapstructure:"localPath" json:"localPath" yaml:"localPath"`
	CloudPath       string `mapstructure:"cloudPath" json:"cloudPath" yaml:"cloudPath"`
	DockerPath      string `mapstructure:"dockerPath" json:"dockerPath" yaml:"dockerPath"`
	InnerDockerPath string `mapstructure:"innerDockerPath" json:"innerDockerPath" yaml:"innerDockerPath"`
	LocalPwd        string `mapstructure:"localPwd" json:"localPwd" yaml:"localPwd"`
	DockerPwd       string `mapstructure:"dockerPwd" json:"dockerPwd" yaml:"dockerPwd"`
}

// 获取dsn配置
func (m *Mysql) DSN(mode string) string {
	if mode == "local" {
		m.Path = m.LocalPath
		m.Password = m.LocalPwd
	} else if mode == "docker" {
		m.Path = m.DockerPath
		m.Password = m.DockerPwd
	} else if mode == "cloud" {
		m.Path = m.CloudPath
		m.Password = m.DockerPwd
	} else if mode == "real" {
		m.Path = m.InnerDockerPath
		m.Password = m.DockerPwd
	}
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
}