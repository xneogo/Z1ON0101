/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/9/30 -- 18:06
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: setting.go
*/

package xmanager

import (
	"fmt"
	"github.com/qiguanzhu/infra/pkg/consts"
	"github.com/qiguanzhu/infra/seele/zsql"
	"net/url"
	"strconv"
	"time"
)

func boolSetting(source, param string, ok bool) string {
	return fmt.Sprintf(consts.CDSNFormat, source, param, strconv.FormatBool(ok))
}

func timeSetting(source, param string, t time.Duration) string {
	// make sure 1ms<=t<24h
	if t < time.Millisecond || t >= 24*time.Hour {
		return ""
	}
	return fmt.Sprintf(consts.CDSNFormat, source, param, t)
}

func stringSetting(source, param, value string) string {
	if "" == value {
		return ""
	}
	return fmt.Sprintf(consts.CDSNFormat, source, param, value)
}

// SetCharset Sets the charset used for client-server interaction
func SetCharset(v string) zsql.Setting {
	return func(source string) string {
		return stringSetting(source, "charset", v)
	}
}

// SetLoc Sets the location for time.Time values (when using parseTime=true). "Local" sets the system's location. See time.LoadLocation for details.
func SetLoc(v string) zsql.Setting {
	return func(source string) string {
		return stringSetting(source, "loc", v)
	}
}

// SetCollation Sets the collation used for client-server interaction on connection. In contrast to charset, collation does not issue additional queries. If the specified collation is unavailable on the target server, the connection will fail.
func SetCollation(v string) zsql.Setting {
	return func(source string) string {
		return stringSetting(source, "collation", v)
	}
}

// SetAllowCleartextPasswords allowCleartextPasswords=true allows using the cleartext client side plugin if required by an account, such as one defined with the PAM authentication plugin. Sending passwords in clear text may be a security problem in some configurations.
func SetAllowCleartextPasswords(ok bool) zsql.Setting {
	return func(source string) string {
		return boolSetting(source, "allowCleartextPasswords", ok)
	}
}

// SetAllowNativePasswords allows the usage of the mysql native password method
func SetAllowNativePasswords(ok bool) zsql.Setting {
	return func(source string) string {
		return boolSetting(source, "allowNativePasswords", ok)
	}
}

// SetAutoCommit set it to true if you know what you are doing
func SetAutoCommit(ok bool) zsql.Setting {
	return func(source string) string {
		return boolSetting(source, "autocommit", ok)
	}
}

// SetParseTime parseTime=true changes the output type of DATE and DATETIME values to time.Time instead of []byte / string
func SetParseTime(ok bool) zsql.Setting {
	return func(source string) string {
		return boolSetting(source, "parseTime", ok)
	}
}

// SetAllowAllFiles allowAllFiles=true disables the file Whitelist for LOAD DATA LOCAL INFILE and allows all files. Might be insecure!
func SetAllowAllFiles(ok bool) zsql.Setting {
	return func(source string) string {
		return boolSetting(source, "allowAllFiles", ok)
	}
}

// SetClientFoundRows clientFoundRows=true causes an UPDATE to return the number of matching rows instead of the number of rows changed.
func SetClientFoundRows(ok bool) zsql.Setting {
	return func(source string) string {
		return boolSetting(source, "clientFoundRows", ok)
	}
}

// SetColumnsWithAlias When columnsWithAlias is true, calls to sql.Rows.Columns() will return the table alias and the column name separated by a dot.
func SetColumnsWithAlias(ok bool) zsql.Setting {
	return func(source string) string {
		return boolSetting(source, "columnsWithAlias", ok)
	}
}

// SetInterpolateParams If interpolateParams is true, placeholders (?) in calls to db.Query() and db.Exec() are interpolated into a single query string with given parameters. This reduces the number of roundtrips, since the driver has to prepare a statement, execute it with given parameters and close the statement again with interpolateParams=false.
func SetInterpolateParams(ok bool) zsql.Setting {
	return func(source string) string {
		return boolSetting(source, "interpolateParams", ok)
	}
}

// SetStrict strict=true enables the strict mode in which MySQL warnings are treated as errors.
func SetStrict(ok bool) zsql.Setting {
	return func(source string) string {
		return boolSetting(source, "strict", ok)
	}
}

// SetTimeout Driver side connection timeout. The value must be a decimal number with an unit suffix ( "ms", "s", "m", "h" ), such as "30s", "0.5m" or "1m30s". To set a server side timeout, use the parameter wait_timeout.
func SetTimeout(timeout time.Duration) zsql.Setting {
	return func(source string) string {
		return timeSetting(source, "timeout", timeout)
	}
}

// SetReadTimeout I/O read timeout. The value must be a decimal number with a unit suffix ( "ms", "s", "m", "h" ), such as "30s", "0.5m" or "1m30s".
func SetReadTimeout(timeout time.Duration) zsql.Setting {
	return func(source string) string {
		return timeSetting(source, "readTimeout", timeout)
	}
}

// SetWriteTimeout I/O write timeout. The value must be a decimal number with a unit suffix ( "ms", "s", "m", "h" ), such as "30s", "0.5m" or "1m30s".
func SetWriteTimeout(timeout time.Duration) zsql.Setting {
	return func(source string) string {
		return timeSetting(source, "writeTimeout", timeout)
	}
}

func GetSettingFunctionList(dynamicConfigure *zsql.DynamicConf) []zsql.Setting {
	return []zsql.Setting{
		SetCharset("utf8mb4"),
		SetCharset("utf8mb4"),
		SetCollation("utf8mb4_unicode_ci"),
		SetAllowCleartextPasswords(true),
		SetInterpolateParams(true),
		SetParseTime(true),
		SetLoc(url.QueryEscape("Asia/Shanghai")),
		SetLoc(url.QueryEscape("Asia/Shanghai")),
		SetTimeout(dynamicConfigure.Timeout),
		SetReadTimeout(dynamicConfigure.ReadTimeout),
		SetWriteTimeout(dynamicConfigure.WriteTimeout),
	}
}
