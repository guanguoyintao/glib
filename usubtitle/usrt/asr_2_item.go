package usrt

import (
	"context"
	ulcs "git.umu.work/AI/uglib/ualgorithm/lcs"
	"git.umu.work/AI/uglib/uregexp"
	"git.umu.work/AI/uglib/usubtitle"
	"git.umu.work/AI/uglib/utypes/ustrings"
	asrpolyproto "git.umu.work/eng/proto/tech/ai/asr_poly/go"
	eastasianwidth "github.com/moznion/go-unicode-east-asian-width"
	"math"
	"time"
)

type ASR2Subtitle struct {
	lineNo        int
	items         []*usubtitle.Item
	minWidth      int
	maxWidth      int
	wordList      []*asrpolyproto.WordInfo
	segmentations []string
}

func NewASR2Subtitle(ctx context.Context, minWidth, maxWidth int, wordList []*asrpolyproto.WordInfo) *ASR2Subtitle {

	return &ASR2Subtitle{
		items:         make([]*usubtitle.Item, 0),
		minWidth:      minWidth,
		maxWidth:      maxWidth,
		wordList:      wordList,
		segmentations: make([]string, 0),
		lineNo:        0,
	}
}

func (s *ASR2Subtitle) isInRange(num, minNum, maxNum int) bool {

	return num >= minNum && num <= maxNum
}

func (s *ASR2Subtitle) sentence2segmentations(ctx context.Context, sentences []string) error {
	for _, sentence := range sentences {
		punctuationBasedSegmentationTexts := ustrings.PunctuationBasedSegmentation(ctx, sentence)
		s.segmentations = append(s.segmentations, punctuationBasedSegmentationTexts...)
	}

	return nil
}

func (s *ASR2Subtitle) constructiveNextSegmentation(ctx context.Context, nextSegmentation string, segmentations []string) string {
	nextSegmentation = uregexp.RemoveSymbolsAndNumbers(nextSegmentation)
	if segmentations == nil || len(segmentations) == 0 {
		return nextSegmentation
	}
	for s.getSentenceSegmentationWidth(ctx, nextSegmentation) <= 20 {
		if len(segmentations) > 0 {
			segmentation := uregexp.RemoveSymbolsAndNumbers(segmentations[0])
			nextSegmentation += segmentation
			if len(segmentations) > 1 {
				segmentations = segmentations[1:]
			} else {
				segmentations = make([]string, 0)
			}
		} else {
			break
		}
	}
	if s.getSentenceSegmentationWidth(ctx, nextSegmentation) > 20 {
		nextSegmentation = string([]rune(nextSegmentation)[:20])
	}

	return nextSegmentation
}

func (s *ASR2Subtitle) getSentenceSegmentationWidth(ctx context.Context, segmentation string) int {
	segmentationSqu := []rune(segmentation)
	var width int
	for _, character := range segmentationSqu {
		if eastasianwidth.IsFullwidth(character) {
			width += 2
		} else {
			width += 1
		}
	}

	return width
}

func (s *ASR2Subtitle) splitSentenceSegmentation(ctx context.Context, segmentation string, splitWidth int) (string, string) {
	segmentationSqu := []rune(segmentation)
	var firstSegmentation string
	var secondSegmentation string
	var width int
	for _, character := range segmentationSqu {
		if eastasianwidth.IsFullwidth(character) {
			width += 2
		} else {
			width += 1
		}
		if width <= splitWidth {
			firstSegmentation += string(character)
		} else {
			secondSegmentation += string(character)
		}
	}

	return firstSegmentation, secondSegmentation
}

func (s *ASR2Subtitle) sentenceSegmentation2Items(ctx context.Context) error {
	segmentation := s.segmentations[0]
	var nextSegmentation string
	if len(s.segmentations) > 1 {
		nextSegmentation = s.constructiveNextSegmentation(ctx, s.segmentations[1], s.segmentations[2:])
		s.segmentations = s.segmentations[1:]
	}
	var remainingLongSegmentation string
	for {
		if len(segmentation) == 0 {
			break
		}
		length := s.getSentenceSegmentationWidth(ctx, segmentation)
		if s.isInRange(length, s.minWidth, 2*s.maxWidth) || len(s.segmentations) == 0 || len(remainingLongSegmentation) > 0 {
			// 构建item
			s.lineNo += 1
			sentenceWordList := s.splitWordListOnSentenceSegmentation(ctx, segmentation, nextSegmentation)
			var item *usubtitle.Item
			if s.isInRange(length, 0, s.maxWidth) {
				item = &usubtitle.Item{
					No:      s.lineNo,
					EndAt:   time.Duration(sentenceWordList[len(sentenceWordList)-1].End) * time.Millisecond,
					StartAt: time.Duration(sentenceWordList[0].Start) * time.Millisecond,
					Texts:   []string{segmentation},
				}
			} else {
				firstSegmentation, secondSegmentation := s.splitSentenceSegmentation(ctx, segmentation, s.maxWidth)
				item = &usubtitle.Item{
					No:      s.lineNo,
					EndAt:   time.Duration(sentenceWordList[len(sentenceWordList)-1].End) * time.Millisecond,
					StartAt: time.Duration(sentenceWordList[0].Start) * time.Millisecond,
					Texts:   []string{firstSegmentation, secondSegmentation},
				}
			}
			s.items = append(s.items, item)
			// 如果还有剩余长文本未处理，接着处理剩下的长文本
			if len(remainingLongSegmentation) > 0 {
				remainingSegmentationLength := s.getSentenceSegmentationWidth(ctx, remainingLongSegmentation)
				if remainingSegmentationLength > 2*s.maxWidth {
					firstSegmentation, secondSegmentation := s.splitSentenceSegmentation(ctx, remainingLongSegmentation, 2*s.maxWidth)
					segmentation = firstSegmentation
					remainingLongSegmentation = secondSegmentation
					nextSegmentation = s.constructiveNextSegmentation(ctx, remainingLongSegmentation, s.segmentations)
				} else {
					segmentation = remainingLongSegmentation
					remainingLongSegmentation = ""
					if len(s.segmentations) == 0 {
						segmentation = ""
						nextSegmentation = ""
					} else {
						nextSegmentation = s.constructiveNextSegmentation(ctx, s.segmentations[0], s.segmentations[1:])
					}
				}
			} else {
				if len(s.segmentations) == 0 {
					segmentation = ""
					continue
				}
				segmentation = s.segmentations[0]
				if len(s.segmentations) > 1 {
					nextSegmentation = s.constructiveNextSegmentation(ctx, s.segmentations[1], s.segmentations[2:])
				} else {
					nextSegmentation = ""
				}
				s.segmentations = s.segmentations[1:]
			}
		} else if s.isInRange(length, 0, s.minWidth) {
			if len(s.segmentations) == 0 {
				continue
			}
			segmentation += s.segmentations[0]
			if len(s.segmentations) > 1 {
				nextSegmentation = s.constructiveNextSegmentation(ctx, s.segmentations[1], s.segmentations[2:])
			} else {
				nextSegmentation = ""
			}
			s.segmentations = s.segmentations[1:]
		} else if s.isInRange(length, 2*s.maxWidth, math.MaxInt) {
			segmentation, remainingLongSegmentation = s.reduceSegmentation(ctx, segmentation)
			nextSegmentation = s.constructiveNextSegmentation(ctx, remainingLongSegmentation, s.segmentations)
		}
	}

	return nil
}

func (s *ASR2Subtitle) reduceSegmentation(ctx context.Context, segmentation string) (string, string) {
	var text string
	for _, word := range s.wordList {
		tmpText := text + word.Word
		length := s.getSentenceSegmentationWidth(ctx, tmpText)
		if length >= s.maxWidth {
			break
		}
		text = tmpText
	}
	length := s.getSentenceSegmentationWidth(ctx, text)
	remainingLongSegmentation := string([]rune(segmentation)[length:])
	splitSegmentation := string([]rune(segmentation)[:length])

	return splitSegmentation, remainingLongSegmentation
}

func (s *ASR2Subtitle) splitWordListOnSentenceSegmentation(ctx context.Context, segmentation, nextSegmentation string) []*asrpolyproto.WordInfo {
	// 最后一句
	var wordIndex int
	if len(nextSegmentation) == 0 {
		wordIndex = len(s.wordList) - 1
	} else {
		// 去除标点
		segmentationWithoutPunctuation := uregexp.RemoveSymbolsAndNumbers(segmentation)
		nextSegmentationWithoutPunctuation := uregexp.RemoveSymbolsAndNumbers(nextSegmentation)
		wordIndex = s.segmentationSimWordIndex(ctx, segmentationWithoutPunctuation, nextSegmentationWithoutPunctuation)
		if wordIndex == -1 {
			return s.wordList
		}
	}
	// 得到两个切分后的列表
	sentenceWordList := s.wordList[:wordIndex+1]
	s.wordList = s.wordList[wordIndex+1:]

	return sentenceWordList
}

// 查找word info list匹配的index
func (s *ASR2Subtitle) segmentationSimWordIndex(ctx context.Context, segmentation, nextSegmentation string) int {
	var activitySegmentation string
	activityWordList := make([]*asrpolyproto.WordInfo, 0)
	maxLcsLength := 0
	for _, wordInfo := range s.wordList {
		activityWord := uregexp.RemoveSymbolsAndNumbers(wordInfo.Word)
		activitySegmentation += activityWord
		activityWordList = append(activityWordList, wordInfo)
		_, nextSegmentationLcsLength := ulcs.Lcs(activitySegmentation, nextSegmentation)
		maxLcsLength = nextSegmentationLcsLength
		if float64(nextSegmentationLcsLength)/float64(s.getSentenceSegmentationWidth(ctx, nextSegmentation)) >= 0.5 {
			break
		}
	}
	activitySegmentation = ""
	// 向前遍历找最大匹配子串，并裁剪activity word list
	removeWordLength := 0
	for i := len(activityWordList) - 1; i >= 0; i-- {
		removeWordLength += 1
		activityWord := uregexp.RemoveSymbolsAndNumbers(activityWordList[i].Word)
		activitySegmentation = activityWord + activitySegmentation
		_, nextSegmentationLcsLength := ulcs.Lcs(nextSegmentation, activitySegmentation)
		if nextSegmentationLcsLength >= maxLcsLength {
			break
		}
	}

	return len(activityWordList) - removeWordLength - 1
}

func (s *ASR2Subtitle) ToItems(ctx context.Context, sentences []string) ([]*usubtitle.Item, error) {
	// 按照标点符号分句
	err := s.sentence2segmentations(ctx, sentences)
	if err != nil {
		return nil, err
	}
	// 通过asr的word info list转换成srt item list
	err = s.sentenceSegmentation2Items(ctx)
	if err != nil {
		return nil, err
	}

	return s.items, nil
}
