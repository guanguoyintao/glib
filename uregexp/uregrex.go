package uregexp

import "regexp"

func match(text string, regex *regexp.Regexp) []string {
	parsed := regex.FindAllString(text, -1)
	return parsed
}

// Date 找到所有日期字符串
func Date(text string) []string {
	return match(text, DateRegex)
}

// Time 找到所有时间字符串
func Time(text string) []string {
	return match(text, TimeRegex)
}

// Phones 找到所有电话号码
func Phones(text string) []string {
	return match(text, PhoneRegex)
}

// PhonesWithExts 找到所有带扩展号的电话号码
func PhonesWithExts(text string) []string {
	return match(text, PhonesWithExtsRegex)
}

// Links 找到所有链接字符串
func Links(text string) []string {
	return match(text, LinkRegex)
}

// Emails 找到所有电子邮件字符串
func Emails(text string) []string {
	return match(text, EmailRegex)
}

// IPv4s 找到所有IPv4地址
func IPv4s(text string) []string {
	return match(text, IPv4Regex)
}

// IPv6s 找到所有IPv6地址
func IPv6s(text string) []string {
	return match(text, IPv6Regex)
}

// IPs 找到所有IP地址（包括IPv4和IPv6）
func IPs(text string) []string {
	return match(text, IPRegex)
}

// NotKnownPorts 找到所有未知端口号
func NotKnownPorts(text string) []string {
	return match(text, NotKnownPortRegex)
}

// Prices 找到所有价格字符串
func Prices(text string) []string {
	return match(text, PriceRegex)
}

// HexColors 找到所有十六进制颜色值
func HexColors(text string) []string {
	return match(text, HexColorRegex)
}

// CreditCards 找到所有信用卡号码
func CreditCards(text string) []string {
	return match(text, CreditCardRegex)
}

// BtcAddresses 找到所有比特币地址
func BtcAddresses(text string) []string {
	return match(text, BtcAddressRegex)
}

// StreetAddresses 找到所有街道地址
func StreetAddresses(text string) []string {
	return match(text, StreetAddressRegex)
}

// ZipCodes 找到所有邮政编码
func ZipCodes(text string) []string {
	return match(text, ZipCodeRegex)
}

// PoBoxes 找到所有邮政信箱字符串
func PoBoxes(text string) []string {
	return match(text, PoBoxRegex)
}

// SSNs 找到所有社会安全号码字符串
func SSNs(text string) []string {
	return match(text, SSNRegex)
}

// MD5Hexes 找到所有MD5十六进制字符串
func MD5Hexes(text string) []string {
	return match(text, MD5HexRegex)
}

// SHA1Hexes 找到所有SHA1十六进制字符串
func SHA1Hexes(text string) []string {
	return match(text, SHA1HexRegex)
}

// SHA256Hexes 找到所有SHA256十六进制字符串
func SHA256Hexes(text string) []string {
	return match(text, SHA256HexRegex)
}

// GUIDs 找到所有GUID字符串
func GUIDs(text string) []string {
	return match(text, GUIDRegex)
}

// ISBN13s 找到所有ISBN13字符串
func ISBN13s(text string) []string {
	return match(text, ISBN13Regex)
}

// ISBN10s 找到所有ISBN10字符串
func ISBN10s(text string) []string {
	return match(text, ISBN10Regex)
}

// VISACreditCards 找到所有VISA信用卡号码
func VISACreditCards(text string) []string {
	return match(text, VISACreditCardRegex)
}

// MCCreditCards 找到所有MasterCard信用卡号码
func MCCreditCards(text string) []string {
	return match(text, MCCreditCardRegex)
}

// MACAddresses 找到所有MAC地址
func MACAddresses(text string) []string {
	return match(text, MACAddressRegex)
}

// IBANs 找到所有IBAN字符串
func IBANs(text string) []string {
	return match(text, IBANRegex)
}

// GitRepos 找到所有带协议前缀的git仓库地址
func GitRepos(text string) []string {
	return match(text, GitRepoRegex)
}

// Punctuations 匹配文本中的标点符号
func Punctuations(text string) []string {
	return match(text, PunctuationRegex)
}

// RemovePunctuation 删除文本中标点
func RemovePunctuation(text string) string {
	re := PunctuationRegex
	return re.ReplaceAllString(text, "")
}

// RemoveSymbolsAndNumbers 删除文本中特殊符号和阿拉伯数字
func RemoveSymbolsAndNumbers(text string) string {
	re := SymbolsAndNumbersRegex
	return re.ReplaceAllString(text, "")
}

// SentenceEndPunctuation 正则表达式用于匹配句子末尾的标点符号
func SentenceEndPunctuation(text string) bool {
	re := SentenceEndPunctuationRegex
	return re.MatchString(text)
}
