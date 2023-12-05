package ulanguage

import (
	"git.umu.work/AI/uglib/ujson"
)

type Language struct {
	BCP47  string `json:"BCP-47"`
	ISO639 string `json:"ISO-639"`
}

func (l *Language) String() string {
	s, _ := ujson.MarshalString(l)

	return s
}

var (
	LanguageInvalid Language = Language{BCP47: "", ISO639: ""}
	// AFZA 南非荷兰语（南非）
	AFZA Language = Language{BCP47: "af-ZA", ISO639: "af"}
	// SQAL 阿尔巴尼亚语（阿尔巴尼亚）
	SQAL Language = Language{BCP47: "sq-AL", ISO639: "sq"}
	// AMET 阿姆哈拉语（埃塞俄比亚）
	AMET Language = Language{BCP47: "am-ET", ISO639: "am"}
	// ARDZ 阿拉伯语（阿尔及利亚）
	ARDZ Language = Language{BCP47: "ar-DZ", ISO639: "ar"}
	// ARIL 阿拉伯语（以色列）
	ARIL Language = Language{BCP47: "ar-IL", ISO639: "ar"}
	// ARJO 阿拉伯语（约旦）
	ARJO Language = Language{BCP47: "ar-JO", ISO639: "ar"}
	// ARKW 阿拉伯语（科威特）
	ARKW Language = Language{BCP47: "ar-KW", ISO639: "ar"}
	// ARLB 阿拉伯语（黎巴嫩）
	ARLB Language = Language{BCP47: "ar-LB", ISO639: "ar"}
	// ARMR 阿拉伯语（毛里塔尼亚）
	ARMR Language = Language{BCP47: "ar-MR", ISO639: "ar"}
	// ARMA 阿拉伯语（摩洛哥）
	ARMA Language = Language{BCP47: "ar-MA", ISO639: "ar"}
	// AROM 阿拉伯语（阿曼）
	AROM Language = Language{BCP47: "ar-OM", ISO639: "ar"}
	// ARQA 阿拉伯语（卡塔尔）
	ARQA Language = Language{BCP47: "ar-QA", ISO639: "ar"}
	// ARSA 阿拉伯语（沙特阿拉伯）
	ARSA Language = Language{BCP47: "ar-SA", ISO639: "ar"}
	// ARPS 阿拉伯语（巴勒斯坦国）
	ARPS Language = Language{BCP47: "ar-PS", ISO639: "ar"}
	// ARTN 阿拉伯语（突尼斯）
	ARTN Language = Language{BCP47: "ar-TN", ISO639: "ar"}
	// AREA 阿拉伯语（阿拉伯联合酋长国）
	AREA Language = Language{BCP47: "ar-AE", ISO639: "ar"}
	// AYEM 阿拉伯语（也门）
	AYEM Language = Language{BCP47: "ar-YE", ISO639: "ar"}
	// HYAM 亚美尼亚语（亚美尼亚）
	HYAM Language = Language{BCP47: "hy-AM", ISO639: "hy"}
	// AZAZ 阿塞拜疆语（阿塞拜疆）
	AZAZ Language = Language{BCP47: "az-AZ", ISO639: "az"}
	// EUES 巴斯克语（西班牙）
	EUES Language = Language{BCP47: "eu-ES", ISO639: "eu"}
	// BNBD 孟加拉语（孟加拉）
	BNBD Language = Language{BCP47: "bn-BD", ISO639: "bn"}
	// BNIN 孟加拉语（印度）
	BNIN Language = Language{BCP47: "bn-IN", ISO639: "bn"}
	// BSBA 波斯尼亚语（波斯尼亚和黑塞哥维那）
	BSBA Language = Language{BCP47: "bs-BA", ISO639: "bs"}
	// BGBG 保加利亚语（保加利亚）
	BGBG Language = Language{BCP47: "bg-BG", ISO639: "bg"}
	// MYMM 缅甸语（缅甸）
	MYMM Language = Language{BCP47: "my-MM", ISO639: "my"}
	// CAES 加泰罗尼亚语（西班牙）
	CAES Language = Language{BCP47: "ca-ES", ISO639: "ca"}
	// ZHHK 中文粤语（香港繁体）
	ZHHK Language = Language{BCP47: "yue-Hant-HK", ISO639: "zh"}
	// ZHCN 中文普通话（中国简体）
	ZHCN Language = Language{BCP47: "zh-Hans-CN", ISO639: "zh"}
	// ZHTW 中文普通话（台湾繁体）
	ZHTW Language = Language{BCP47: "zh-Hant-TW", ISO639: "zh"}
	// HRHR 克罗地亚语（克罗地亚）
	HRHR Language = Language{BCP47: "hr-HR", ISO639: "hr"}
	// CSCZ 捷克语（捷克共和国）
	CSCZ Language = Language{BCP47: "cs-CZ", ISO639: "cs"}
	// DADK 丹麦语（丹麦）
	DADK Language = Language{BCP47: "da-DK", ISO639: "da"}
	// NLBE 荷兰语（比利时）
	NLBE Language = Language{BCP47: "nl-BE", ISO639: "nl"}
	// NLNL 荷兰语（荷兰）
	NLNL Language = Language{BCP47: "nl-NL", ISO639: "nl"}
	// ENAU 英语（澳大利亚）
	ENAU Language = Language{BCP47: "en-AU", ISO639: "en"}
	// ENCA 英语（加拿大）
	ENCA Language = Language{BCP47: "en-CA", ISO639: "en"}
	// ENGH 英语（加纳）
	ENGH Language = Language{BCP47: "en-GH", ISO639: "en"}
	// ENHK 英语（香港）
	ENHK Language = Language{BCP47: "en-HK", ISO639: "en"}
	// ENIN 英语（印度）
	ENIN Language = Language{BCP47: "en-IN", ISO639: "en"}
	// ENIE 英语（爱尔兰）
	ENIE Language = Language{BCP47: "en-IE", ISO639: "en"}
	// ENKE 英语（肯尼亚）
	ENKE Language = Language{BCP47: "en-KE", ISO639: "en"}
	// ENNZ 英语（新西兰）
	ENNZ Language = Language{BCP47: "en-NZ", ISO639: "en"}
	// ENNG 英语（尼日利亚）
	ENNG Language = Language{BCP47: "en-NG", ISO639: "en"}
	// ENPK 英语（巴基斯坦）
	ENPK Language = Language{BCP47: "en-PK", ISO639: "en"}
	// ENPH 英语（菲律宾）
	ENPH Language = Language{BCP47: "en-PH", ISO639: "en"}
	// ENSG 英语（新加坡）
	ENSG Language = Language{BCP47: "en-SG", ISO639: "en"}
	// ENZA 英语（南非）
	ENZA Language = Language{BCP47: "en-ZA", ISO639: "en"}
	// ENTZ 英语（坦桑尼亚）
	ENTZ Language = Language{BCP47: "en-TZ", ISO639: "en"}
	// ENGB 英语（英国）
	ENGB Language = Language{BCP47: "en-GB", ISO639: "en"}
	// ENUS 英语（美国）
	ENUS Language = Language{BCP47: "en-US", ISO639: "en"}
	// ETET 爱沙尼亚语（爱沙尼亚）
	ETET Language = Language{BCP47: "et-EE", ISO639: "et"}
	// FILPH 菲律宾语（菲律宾）
	FILPH Language = Language{BCP47: "fil-PH", ISO639: "fil"}
	// FIFI 芬兰语（芬兰）
	FIFI Language = Language{BCP47: "fi-FI", ISO639: "fi"}
	// FRBE 法语（比利时）
	FRBE Language = Language{BCP47: "fr-BE", ISO639: "fr"}
	// FRCA 法语（加拿大）
	FRCA Language = Language{BCP47: "fr-CA", ISO639: "fr"}
	// FRFR 法语（法国）
	FRFR Language = Language{BCP47: "fr-FR", ISO639: "fr"}
	// FRCH 法语（瑞士）
	FRCH Language = Language{BCP47: "fr-CH", ISO639: "fr"}
	// GLES 加利西亚语（西班牙）
	GLES Language = Language{BCP47: "gl-ES", ISO639: "gl"}
	// KAGE 格鲁吉亚语（格鲁吉亚）
	KAGE Language = Language{BCP47: "ka-GE", ISO639: "ka"}
	// DEAT 德语（奥地利）
	DEAT Language = Language{BCP47: "de-AT", ISO639: "de"}
	// DEDE 德语（德国）
	DEDE Language = Language{BCP47: "de-DE", ISO639: "de"}
	// DECH 德语（瑞士）
	DECH Language = Language{BCP47: "de-CH", ISO639: "de"}
	// ELGR 希腊语（希腊）
	ELGR Language = Language{BCP47: "el-GR", ISO639: "el"}
	// GUIN 古吉拉特语（印度）
	GUIN Language = Language{BCP47: "gu-IN", ISO639: "gu"}
	// IWIL 希伯来语（以色列）
	IWIL Language = Language{BCP47: "iw-IL", ISO639: "iw"}
	// HIIN 印地语（印度）
	HIIN Language = Language{BCP47: "hi-IN", ISO639: "hi"}
	// HUHU 匈牙利语（匈牙利）
	HUHU Language = Language{BCP47: "hu-HU", ISO639: "hu"}
	// ISIS 冰岛语（冰岛）
	ISIS Language = Language{BCP47: "is-IS", ISO639: "is"}
	// IDID 印度尼西亚语（印度尼西亚）
	IDID Language = Language{BCP47: "id-ID", ISO639: "id"}
	// ITIT 意大利语（意大利）
	ITIT Language = Language{BCP47: "it-IT", ISO639: "it"}
	// ITCH 意大利语（瑞士）
	ITCH Language = Language{BCP47: "it-CH", ISO639: "it"}
	// JAJP 日语（日本）
	JAJP Language = Language{BCP47: "ja-JP", ISO639: "ja"}
	// JVID 爪哇语（印度尼西亚）
	JVID Language = Language{BCP47: "jv-ID", ISO639: "jv"}
	// KNIN 卡纳达语（印度）
	KNIN Language = Language{BCP47: "kn-IN", ISO639: "kn"}
	// KKKZ 哈萨克语（哈萨克斯坦）
	KKKZ Language = Language{BCP47: "kk-KZ", ISO639: "kk"}
	// KMKH 高棉语（柬埔寨）
	KMKH Language = Language{BCP47: "km-KH", ISO639: "km"}
	// RWRW 卢旺达语（卢旺达）
	RWRW Language = Language{BCP47: "rw-RW", ISO639: "rw"}
	// KOKR 韩语（韩国）
	KOKR Language = Language{BCP47: "ko-KR", ISO639: "ko"}
	// LOLA 老挝语（老挝）
	LOLA Language = Language{BCP47: "lo-LA", ISO639: "lo"}
	// LVLV 拉脱维亚语（拉脱维亚）
	LVLV Language = Language{BCP47: "lv-LV", ISO639: "lv"}
	// LTLT 立陶宛语（立陶宛）
	LTLT Language = Language{BCP47: "lt-LT", ISO639: "lt"}
	// MKMK 马其顿语（北马其顿）
	MKMK Language = Language{BCP47: "mk-MK", ISO639: "mk"}
	// MSMY 马来语（马来西亚）
	MSMY Language = Language{BCP47: "ms-MY", ISO639: "ms"}
	// MLIN 马拉雅拉姆语（印度）
	MLIN Language = Language{BCP47: "ml-IN", ISO639: "ml"}
	// MRIN 马拉地语（印度）
	MRIN Language = Language{BCP47: "mr-IN", ISO639: "mr"}
	// MNMN 蒙古语（蒙古）
	MNMN Language = Language{BCP47: "mn-MN", ISO639: "mn"}
	// NENP 尼泊尔语（尼泊尔）
	NENP Language = Language{BCP47: "ne-NP", ISO639: "ne"}
	// NONO 博克马尔挪威语（挪威）
	NONO Language = Language{BCP47: "no-NO", ISO639: "no"}
	// FAIR 波斯语（伊朗）
	FAIR Language = Language{BCP47: "fa-IR", ISO639: "fa"}
	// PLPL 波兰语（波兰）
	PLPL Language = Language{BCP47: "pl-PL", ISO639: "pl"}
	// PTBR 葡萄牙语（巴西）
	PTBR Language = Language{BCP47: "pt-BR", ISO639: "pt"}
	// PTPT 葡萄牙语（葡萄牙）
	PTPT Language = Language{BCP47: "pt-PT", ISO639: "pt"}
	// PAGURU 旁遮普语（果鲁穆奇语，印度）
	PAGURU Language = Language{BCP47: "pa-Guru-IN", ISO639: "pa"}
	// RORO 罗马尼亚语（罗马尼亚）
	RORO Language = Language{BCP47: "ro-RO", ISO639: "ro"}
	// RURU 俄语（俄罗斯）
	RURU Language = Language{BCP47: "ru-RU", ISO639: "ru"}
	// SRSR 塞尔维亚语（塞尔维亚）
	SRSR Language = Language{BCP47: "sr-RS", ISO639: "sr"}
	// SILK 僧伽罗语（斯里兰卡）
	SILK Language = Language{BCP47: "si-LK", ISO639: "si"}
	// SKSK 斯洛伐克语（斯洛伐克）
	SKSK Language = Language{BCP47: "sk-SK", ISO639: "sk"}
	// SLSI 斯洛文尼亚语（斯洛文尼亚）
	SLSI Language = Language{BCP47: "sl-SI", ISO639: "sl"}
	// STZA 南索托语（南非）
	STZA Language = Language{BCP47: "st-ZA", ISO639: "st"}
	// ESAR 西班牙语（阿根廷）
	ESAR Language = Language{BCP47: "es-AR", ISO639: "es"}
	// ESBO 西班牙语（玻利维亚）
	ESBO Language = Language{BCP47: "es-BO", ISO639: "es"}
	// ESCL 西班牙语（智利）
	ESCL Language = Language{BCP47: "es-CL", ISO639: "es"}
	// ESCO 西班牙语（哥伦比亚）
	ESCO Language = Language{BCP47: "es-CO", ISO639: "es"}
	// ESCR 西班牙语（哥斯达黎加）
	ESCR Language = Language{BCP47: "es-CR", ISO639: "es"}
	// ESDO 西班牙语（多米尼加共和国）
	ESDO Language = Language{BCP47: "es-DO", ISO639: "es"}
	// ESEC 西班牙语（厄瓜多尔）
	ESEC Language = Language{BCP47: "es-EC", ISO639: "es"}
	// ESSV 西班牙语（萨尔瓦多）
	ESSV Language = Language{BCP47: "es-SV", ISO639: "es"}
	// ESGT 西班牙语（危地马拉）
	ESGT Language = Language{BCP47: "es-GT", ISO639: "es"}
	// ESHN 西班牙语（洪都拉斯）
	ESHN Language = Language{BCP47: "es-HN", ISO639: "es"}
	// ESMX 西班牙语（墨西哥）
	ESMX Language = Language{BCP47: "es-MX", ISO639: "es"}
	// ESNI 西班牙语（尼加拉瓜）
	ESNI Language = Language{BCP47: "es-NI", ISO639: "es"}
	// ESPA 西班牙语（巴拿马）
	ESPA Language = Language{BCP47: "es-PA", ISO639: "es"}
	// ESPY 西班牙语（巴拉圭）
	ESPY Language = Language{BCP47: "es-PY", ISO639: "es"}
	// ESPE 西班牙语（秘鲁）
	ESPE Language = Language{BCP47: "es-PE", ISO639: "es"}
	// ESPR 西班牙语（波多黎各）
	ESPR Language = Language{BCP47: "es-PR", ISO639: "es"}
	// ESES 西班牙语（西班牙）
	ESES Language = Language{BCP47: "es-ES", ISO639: "es"}
	// ESUS 西班牙语（美国）
	ESUS Language = Language{BCP47: "es-US", ISO639: "es"}
	// ESUY 西班牙语（乌拉圭）
	ESUY Language = Language{BCP47: "es-UY", ISO639: "es"}
	// ESVE 西班牙语（委内瑞拉）
	ESVE Language = Language{BCP47: "es-VE", ISO639: "es"}
	// SUID 巽他语（印度尼西亚）
	SUID Language = Language{BCP47: "su-ID", ISO639: "su"}
	// SWKE 斯瓦希里语（肯尼亚）
	SWKE Language = Language{BCP47: "sw-KE", ISO639: "sw"}
	// SWTZ 斯瓦希里语（坦桑尼亚）
	SWTZ Language = Language{BCP47: "sw-TZ", ISO639: "sw"}
	// SSLATZA 斯威士语（拉丁字母，南非）
	SSLATZA Language = Language{BCP47: "ss-Latn-ZA", ISO639: "ss"}
	// SVSE 瑞典语（瑞典）
	SVSE Language = Language{BCP47: "sv-SE", ISO639: "sv"}
	// TAIND 泰米尔语（印度）
	TAIND Language = Language{BCP47: "ta-IN", ISO639: "ta"}
	// TAMY 泰米尔语（马来西亚）
	TAMY Language = Language{BCP47: "ta-MY", ISO639: "ta"}
	// TASN 泰米尔语（新加坡）
	TASN Language = Language{BCP47: "ta-SG", ISO639: "ta"}
	// TALK 泰米尔语（斯里兰卡）
	TALK Language = Language{BCP47: "ta-LK", ISO639: "ta"}
	// TEIN 泰卢固语（印度）
	TEIN Language = Language{BCP47: "te-IN", ISO639: "te"}
	// THTH 泰语（泰国）
	THTH Language = Language{BCP47: "th-TH", ISO639: "th"}
	// TSZA 聪加语（南非）
	TSZA Language = Language{BCP47: "ts-ZA", ISO639: "ts"}
	// TNLATNZA 茨瓦纳语（拉丁字母，南非）
	TNLATNZA Language = Language{BCP47: "tn-Latn-ZA", ISO639: "tn"}
	// TRTR 土耳其语（土耳其）
	TRTR Language = Language{BCP47: "tr-TR", ISO639: "tr"}
	// UKUA 乌克兰语（乌克兰）
	UKUA Language = Language{BCP47: "uk-UA", ISO639: "uk"}
	// URIN 乌尔都语（印度）
	URIN Language = Language{BCP47: "ur-IN", ISO639: "ur"}
	// URPK 乌尔都语（巴基斯坦）
	URPK Language = Language{BCP47: "ur-PK", ISO639: "ur"}
	// UZUZ 乌兹别克语（乌兹别克斯坦）
	UZUZ Language = Language{BCP47: "uz-UZ", ISO639: "uz"}
	// VEZA 文达语（南非）
	VEZA Language = Language{BCP47: "ve-ZA", ISO639: "ve"}
	// VIVN 越南语（越南）
	VIVN Language = Language{BCP47: "vi-VN", ISO639: "vi"}
	// XHZA 科萨语（南非）
	XHZA Language = Language{BCP47: "xh-ZA", ISO639: "xh"}
	// ZUZA 祖鲁语（南非）
	ZUZA Language = Language{BCP47: "zu-ZA", ISO639: "zu"}
)

var AllLanguages = []Language{
	AFZA, SQAL, AMET, ARDZ, ARIL, ARJO, ARKW, ARLB, ARMR, ARMA, AROM, ARQA, ARSA, ARPS, ARTN, AREA, AYEM,
	HYAM, AZAZ, EUES, BNBD, BNIN, BSBA, BGBG, MYMM, CAES, ZHHK, ZHCN, ZHTW, HRHR, CSCZ, DADK, NLBE, NLNL,
	ENAU, ENCA, ENGH, ENHK, ENIN, ENIE, ENKE, ENNZ, ENNG, ENPK, ENPH, ENSG, ENZA, ENTZ, ENGB, ENUS, ETET,
	FILPH, FIFI, FRBE, FRCA, FRFR, FRCH, GLES, KAGE, DEAT, DEDE, DECH, ELGR, GUIN, IWIL, HIIN, HUHU, ISIS,
	IDID, ITIT, ITCH, JAJP, JVID, KNIN, KKKZ, KMKH, RWRW, KOKR, LOLA, LVLV, LTLT, MKMK, MSMY, MLIN, MRIN,
	MNMN, NENP, NONO, FAIR, PLPL, PTBR, PTPT, PAGURU, RORO, RURU, SRSR, SILK, SKSK, SLSI, STZA, ESAR, ESBO,
	ESCL, ESCO, ESCR, ESDO, ESEC, ESSV, ESGT, ESHN, ESMX, ESNI, ESPA, ESPY, ESPE, ESPR, ESES, ESUS, ESUY,
	ESVE, SUID, SWKE, SWTZ, SSLATZA, SVSE, TAIND, TAMY, TASN, TALK, TEIN, THTH, TSZA, TNLATNZA, TRTR, UKUA,
	URIN, URPK, UZUZ, VEZA, VIVN, XHZA, ZUZA,
}

func GetLanguage(bcp47 string) Language {
	for _, language := range AllLanguages {
		if language.BCP47 == bcp47 {
			return language
		}
	}

	return LanguageInvalid
}
