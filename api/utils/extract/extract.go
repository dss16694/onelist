package extract

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func removeEndingOne(s string) string {
	if len(s) > 0 && s[len(s)-1] == '1' {
		return s[:len(s)-1]
	}
	return s
}

// 过滤电影文件名
func ExtractMovieName(s string) string {
	oldName := s
	// 删除发布年份和文件扩展名
	re := regexp.MustCompile(`\d{4}`)
	s = re.ReplaceAllString(s, "")
	// 兼容全是纯数字的
	if len(s) == 0 {
		re = regexp.MustCompile(`\d+`)
		matchenNumbers := re.FindAllString(oldName, -1)
		if len(matchenNumbers) > 0 {
			s = matchenNumbers[0]
		}
	}
	// 删除括号及其内容
	re = regexp.MustCompile(`\s*\([^)]+\)`)
	s = re.ReplaceAllString(s, "")

	// 提取中文名称
	re = regexp.MustCompile(`[\p{Han}\d{1,2}]+`)
	matches := re.FindAllString(s, -1)
	if len(matches) > 0 {
		name := removeEndingOne(matches[0])
		return name
	}
	return s
}

// 根据文件名获取剧集季及集信息
func ExtractNumberWithFile(file string) (int, int, error) {
	p, err := filepath.Abs(file)
	if err != nil {
		return 0, 0, err
	}
	SeasonNumber := 0
	EpisodeNumber := 0
	fileName := filepath.Base(p)
	re := regexp.MustCompile(`[Ss](\d{1,2})[Ee](\d{1,4})`)
	match := re.FindStringSubmatch(fileName)
	season := ""
	episode := ""
	if len(match) < 3 {
		season, err = ExtractSeason(file)
		if err != nil {
			return 0, 0, errors.New("get number error")
		} else {
			re = regexp.MustCompile(`(\d{1,4})`)
			match = re.FindStringSubmatch(fileName)
			episode = match[1]
		}
	} else {
		season = match[1]
		episode = match[2]
	}
	SeasonNumber, err = strconv.Atoi(season)
	if err != nil {
		return 0, 0, err
	}
	EpisodeNumber, err = strconv.Atoi(episode)
	if err != nil {
		return 0, 0, err
	}
	return SeasonNumber, EpisodeNumber, nil
}

// ExtractSeason 提取路径中的季节信息
func ExtractSeason(filePath string) (string, error) {

	filePath, err := url.QueryUnescape(filePath)
	if err != nil {
		return "", fmt.Errorf("get QueryUnescape error")
	}

	// 定义正则表达式以匹配“第X季”、“SXX”等格式
	re := regexp.MustCompile(`第(\p{Han}+季)`)

	// 查找匹配
	matches := re.FindStringSubmatch(filePath)
	if len(matches) > 0 {
		arabicNum, exists := chineseToArabic[matches[0]]
		if exists {
			return arabicNum, nil
		}
	} else {
		re = regexp.MustCompile(`[Ss](\d{1,2})`)
		matches = re.FindStringSubmatch(filePath)
		if len(matches) == 0 {
			index := strings.Index(filePath, "/SP/")
			if index != -1 {
				return "0", nil
			}
			return "1", nil
		}
		return matches[1], nil
	}

	return "", fmt.Errorf("get ExtractSeason error")
}

// 中文数字到阿拉伯数字的映射
var chineseToArabic = map[string]string{
	"第一季": "1",
	"第二季": "2",
	"第三季": "3",
	"第四季": "4",
	"第五季": "5",
	"第六季": "6",
	"第七季": "7",
	"第八季": "8",
	"第九季": "9",
	"第十季": "10",
}
