package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = true

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})

	t.Run("less than 10 words in the input text", func(t *testing.T) {
		testCases := []struct {
			text     string
			expected []string
		}{
			{
				text: "some text which is not long",
				expected: []string{
					"is",
					"long",
					"not",
					"some",
					"text",
					"which",
				},
			},
			{
				text: "some",
				expected: []string{
					"some",
				},
			},
			{
				text: "some text is not long",
				expected: []string{
					"is",
					"long",
					"not",
					"some",
					"text",
				},
			},
		}
		for _, tc := range testCases {
			require.Equal(t, tc.expected, Top10(tc.text))
		}
	})

	t.Run("more than 10 words, same frequency", func(t *testing.T) {
		testText := `ggg ggg ggg ggg ggg ggg ggg ggg ggg ggg
					fff fff fff fff fff fff fff fff fff fff
					ttt ttt ttt ttt ttt ttt ttt ttt ttt ttt
					rrr rrr rrr rrr rrr rrr rrr rrr rrr rrr
					nnn nnn nnn nnn nnn nnn nnn nnn nnn nnn
					bbb bbb bbb bbb bbb bbb bbb bbb bbb bbb
					ccc ccc ccc ccc ccc ccc ccc ccc ccc ccc
					ppp ppp ppp ppp ppp ppp ppp ppp ppp ppp
					ddd ddd ddd ddd ddd ddd ddd ddd ddd ddd
					sss sss sss sss sss sss sss sss sss sss
					eee eee eee eee eee eee eee eee eee eee
					hhh hhh hhh hhh hhh hhh hhh hhh hhh hhh
					aaa aaa aaa aaa aaa aaa aaa aaa aaa aaa
					iii iii iii iii iii iii iii iii iii iii
					kkk kkk kkk kkk kkk kkk kkk kkk kkk kkk
					lll lll lll lll lll lll lll lll lll lll
					mmm mmm mmm mmm mmm mmm mmm mmm mmm mmm
					ooo ooo ooo ooo ooo ooo ooo ooo ooo ooo
					`
		expected := []string{
			"aaa",
			"bbb",
			"ccc",
			"ddd",
			"eee",
			"fff",
			"ggg",
			"hhh",
			"iii",
			"kkk",
		}
		require.Equal(t, expected, Top10(testText))
	})

	t.Run("consecutive punctuation", func(t *testing.T) {
		testText := `hello I would like to say hello,, but there are two commas`
		expected := []string{
			"hello", // 2 times, others in lexicographical order
			"are",
			"but",
			"commas",
			"i",
			"like",
			"say",
			"there",
			"to",
			"two",
		}
		require.Equal(t, expected, Top10(testText))
	})

	t.Run("quotes", func(t *testing.T) {
		testText := `hello I would like to say "hello",, but there are two commas`
		expected := []string{
			"hello", // 2 times, others in lexicographical order
			"are",
			"but",
			"commas",
			"i",
			"like",
			"say",
			"there",
			"to",
			"two",
		}
		require.Equal(t, expected, Top10(testText))
	})
}
