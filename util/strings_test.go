package util

import (
	"testing"
	"strings"
)

const (
	data = `1. 你吃饭了吗？ Nǐ chīfàn le ma?
Literally: “Have you eaten?”
Function: Expresses one’s concern for someone else.
Near-equivalent phrase in English: “How’s it going?” or “How are you?”

2. 你多吃一点。Nǐ duō chī yīdiǎn.
Literally: “Eat some more.”
Function: Expresses one’s hospitality for a guest.
Near-equivalent phrase in English: “Have some more.”

3. 慢慢吃。Màn man chī.
Literally: “Eat slowly.”
Function: Expresses politeness to someone when eating.
Near-equivalent phrase in English: “Bon appétit” or “enjoy your meal” (American English).

4. 慢走。Màn zǒu.
Literally: “Walk slowly.”
Function: Expresses politeness to someone when they leave someone’s house or a hotel, restaurant, etc.
Near-equivalent phrase in English: “Take care” or “Have a good day” (American English).

5. 慢慢来。Màn màn lái.
Literally: “Come slowly.”
Function: Expresses to someone to take it easy.
Near-equivalent phrase in English: “Take it easy”, “Take your time” or “Easy does it”.

6. 我跟你讲。Wǒ gēn nǐ jiǎng.
Literally: “I speak to you.”
Function: Used to get someone to listen to you when you want to tell them something you think is important.
Near-equivalent phrase in English: “Look, …” or “Listen, …”

7. 我先走了。Wǒ xiān zǒu le.
Literally: “I go first.”
Function: Used to tell someone that you are leaving, and that they can stay in the same place if they wish.
Near-equivalent phrase in English: “I’m off.” or “I gotta run.”

8. 请问一下。Qǐng wèn yīxià.
Literally: “Please [let me] ask.”
Function: Used when you wish to ask someone (usually a stranger) a question.
Near-equivalent phrase in English: “Excuse me.”

9. 别送了。Bié sòng le.
Literally: “Don’t see me out.”
Function: Very polite. The guest says this to the host when the guest feels it’s not necessary for the host to see them out.
Near-equivalent phrase in English: “You don’t need to see me out.” or “No need to walk me out.”

10. 我敬你一杯。Wǒ jìng nǐ yī bēi.
Literally: This phrase is difficult to translate literally. _ here symbolises respect given to the second party.
Function: Said when you wish to raise your drink to someone, to drink with them or propose a toast.
Near-equivalent phrase in English: “I drink to you” or just “Cheers”.

11. 我会考虑一下的。Wǒ huì kǎolǜ yīxià de.
Literally: “I will consider [it].”
Function: Used to let someone know that you’ll think about something they have suggested, especially if you’re not really sure you accept it.
Near-equivalent phrase in English: “I’ll think about it.”

12. 你去忙你的吧。Nǐ qù máng nǐ de ba.
Literally: “You go do what you are busy with.”
Function: Used to let someone know that they can continue doing what they are doing, while you go and do something else.
Near-equivalent phrase in English: “Please carry on with what you’re doing.”

13. 我不是说你。Wǒ bù shì shuō nǐ.
Literally: “I’m not criticising you.”
Function: Used to preface something critical you’re about to say and urge the other person not to be offended by it.
Near-equivalent phrase in English: “I’m not criticising you.” or “I’m not having a go at you.” (Aussie English) or “No offense.”

14. 至于吗？ Zhìyú ma?
Literally: Difficult to translate literally; _ is a verb used to indicate that something has reached a certain level, while _ ma creates a question structure.
Function: Used to express doubt about what someone says. You may reply as _.
Near-equivalent phrase in English: “Is that really the case?” or “Has it come to that?” (depending on situation)

15. 你吓死我了。Nǐ xià sǐ wǒ le.
Literally: “You scared me to death.”
Function: Used to express one’s fear or concern about someone.
Near-equivalent phrase in English: “You scared the crap outta me” or “You freaked me out” or “You made me concerned” depending on situation.

16. 随你了。Suí nǐ le.
Literally: “I sui [follow? go with?] you.”
Function: Used to express that, when it comes to making a particular decision, you don’t really mind either way.
Near-equivalent phrase in English: “Up to you” or “I’m easy.” “Whatever/I don’t care” depending on the situation.

17. 来来来… 坐坐坐… 吃吃吃… Lái lái lái…zuò zuò zuò …chī chī chī…
Literally: “Come come come… sit sit sit… eat eat eat
Function: These three different phrases are used in different situations, though they may be said after one another. They are normally used when greeting a guest and you wish to show them your hospitality – to come in and/or take a seat and/or eat.
Near-equivalent phrase in English: “Make yourself at home… Please, take a seat… Tuck in.”

18. [某人]不在状态。[Somebody] bù zài zhuàngtài.
Literally: “Somebody is not in [a normal] state.”
Function: Used to explain that someone – perhaps a friend or a family member – is not feeling very well.
Near-equivalent phrase in English: “Somebody is not him/herself.”

19. 我失陪了。Wǒ shīpéi le.
Literally: “I lose [your] company.”
Function: Used to politely let someone know that you are leaving.
Near-equivalent phrase in English: “I’m sorry but I must take my leave” (very formal) or “Sorry but I have to run” (informal).

20. 请教一下。Qǐngjiào yīxià.
Literally: “Please instruct [me].”
Function: Used to let someone know that you welcome comments and criticism, particularly about a project you have been working on,  your performance, etc.
Near-equivalent phrase in English: “I’d love to hear some feedback from you.”, “I look forward to hearing your advice.”, “Feel free to leave some comments.” etc.

21. 你辛苦了。Nǐ xīnkǔ le.
Literally: “You’ve tasted bitterness/hardship.”
Function: Used to express gratitude for the help someone has given you.
Near-equivalent phrase in English: No real equivalent in English. The translation “You’ve worked so hard.” is acceptable, but probably sounds a little strange. In this situation an English speaker would probably just say, “Thank you so much, I really appreciate it.”

22. [某人]吃了很多苦。[Somebody] chī le hěn duō kǔ.
Literally: “Somebody has eaten a lot of bitterness (hardship).”
Function: Used to state that someone has gone through many hardships.
Near-equivalent phrase in English: “Somebody‘s been through a lot.” or “Somebody has gone through a rough time.”

23. 我听你的。Wǒ tīng nǐ de.
Literally: “I’ll listen to you.”
Function: Used to express that you will listen and follow what someone does, usually for our own good.
Near-equivalent phrase in English: “You’re the boss.”

24. [某人]都还给老师了。 [Something] dōu huán gěi lǎoshī le.
Literally: “Something has all been given back to the teacher.”
Function: Used to indicate that everything that you’ve learnt has been forgotten.
Near-equivalent phrase in English: As far as I know, no real equivalent. “I’ve forgotten it all” would suffice as a reference translation. A native English speaker may say something like, “My French/mathematics/etc is a bit rusty” though this is not as strong as the original Chinese sentence.

25. A生了B的气。A shēng le B de qì.
Literally: “A generated anger because of B.”
Function: Used to express that you have made somebody angry. Notable because this structure in Mandarin is unusual and a little confusing for Chinese learners.
Near-equivalent phrase in English: “A is angry at B.” or “A is pissed off with B.” or “B made A angry.”

26. [某事]不关[某人]的事。[Something] bù guān [somebody] de shì.
Literally: “Something does not relate to the affairs of somebody.”
Function: Used to (quite rudely) point out that something is not the business of someone else.
Near-equivalent phrase in English: “Something is not someone’s business.”. When used as an interjection the phrases “None of your business!” or “What’s it to you?” come to mind – that’s _ nǐ pì shì? in Mandarin.

27. [某人]真够朋友。[Somebody] zhēn gòu péngyǒu.
Literally: “Somebody is really an adequate friend.”
Function: Used to let someone know that you really value their friendship.
Near-equivalent phrase in English: “Somebody is a true friend” or “Somebody is a real mate” in Aussie English.

28. 话不是这么说。Huà bù shì zhème shuō.
Literally: “It is not said like this.”
Function: Used to gently disagree with someone.
Near-equivalent phrase in English: “I don’t really think that’s the case.”

29. 可不是吗？ Kě bù shì ma?
Literally: “How can it not be?”
Function: Used to express your strong agreement about something.
Near-equivalent phrase in English: “Definitely!” or “Absolutely!”

30. 哪儿跟哪儿？ Nǎr gēn nǎr?
Literally: “Where compared to where?”
Function: Used to express doubt about the relationship of two things which you think are not related.
Near-equivalent phrase in English: “I don’t see the connection” or “What’s that got to do with it?”

31. 真有你的。Zhēn yǒu nǐ de.
Literally: _ (“really”) + _ (“you”) + _ (“your [skill; talent]”)
Function: Used to express your admiration of someone’s skill or talent.
Near-equivalent phrase in English: “You’re really awesome.” or “You’re really something else.”

32. 看情况。Kàn qíngkuàng.
Literally: “Look at the situation.”
Function: Used to express uncertainty about a certain situation.
Near-equivalent phrase in English: “Play it by ear” or “It depends” depending on situation

33. 谁跟谁啊？Shéi gēn shéi a?
Literally: “Who with who ah?”
Function: Used to remind the other person that you are good friends with them, to get them to stop being so polite or to get them to reveal to you something you want to know.
Near-equivalent phrase in English: “Come on, we’re friends aren’t we?”

34. [某事]包在我身上。[Something] bāo zài wǒ shēnshang.
Literally: “Something‘s package is on my person.”
Function: Used to let someone know that you will take absolute responsibility for a certain task.
Near-equivalent phrase in English: “Leave it all to me and I’ll make it happen.”

35. [某人]不是东西。[Somebody] bù shì dōngxi.
Literally: “Somebody is not a thing.”
Function: Used to insult someone.
Near-equivalent phrase in English: “Somebody is good-for-nothing.”

36. 就那么回事。Jiù nàme huí shì.
Literally: “That’s how it was.”
Function: To state that something is mediocre or average.
Near-equivalent phrases in English: “Not that great.” or “Average.”

37. [某人]死的心都有。[Somebody] sǐ de xīn dōu yǒu.
Literally: “Somebody even has a dead heart.” (As if their heart is dead.)
Function: Used to express somebody’s desperation, disappointment and/or grief.
Near-equivalent phrase in English: “Somebody is torn apart.”

38. 爱谁谁！ Ài shéishéi!
Literally: “Love who who!”
Function: Used to express indifference.
Near-equivalent phrase in English: “Whatever!” or “Who cares!”

39. [某人]不好那口。[Somebody] bù hào nà kǒu.
Literally: “Somebody is not well (used to) that mouth.”
Function: Used to express that someone does not share a particular hobby or fondness for something.
Near-equivalent phrase in English: “Somebody is not into that.” or “That’s not somebody’s thing.”

40. 不要放在心上。Bù yào fàng zài xīn shàng.
Literally: “Don’t put [it] in [your] heart.”
Function: Used to advise someone to not continue thinking about an unpleasant topic.
Near-equivalent phrase in English: “Don’t take it to heart.”

41. 请你多多包涵。Qǐng nǐ duōduō bāohan.
Literally: “Please forgive [me] much.”
Function: Said before or after you do or say something which you think may hurt or offend others.
Near-equivalent phrase in English: “Please forgive me.” or “Please bear with me.”

42. 给[某人]点儿颜色看看。Gěi [somebody] diǎnr yánsè kàn kàn.
Literally: “Give somebody a little colour (facial expression) to see.”
Function: Used to express someone’s ferociousness, to intimidate someone, usually to warn them that they are tough and not to be offended.
Near-equivalent phrase in English: “Teach someone a lesson.”

43. [某人]的鼻子气歪了。[Somebody] de bízi qì wāi le.
Literally: “Somebody‘s nose is crooked with anger.”
Function: Used to express how angry someone is.
Near-equivalent phrase in English: “He’s really pissed off.”

44. [关于某事]打一个问号。[About something] dǎ yī gè wènhào.
Literally: “About something [I] write a question mark.”
Function: Used to express doubt about something.
Near-equivalent phrase in English: Not sure of an idiomatic equivalent; a basic translation is “to be unsure about something.”

45. [某人]也有今天。[Somebody] yě yǒu jīntiān.
Literally: “Somebody also has today.”
Function: Used to state that someone has gotten comeuppance for a wrong deed.
Near-equivalent phrase in English: “Somebody will get his/her just deserts.” or “Somebody has got what he/she deserves.”`
)

func TestContainsCJK(t *testing.T) {
	lines := strings.Split(data, "\n")
	for idx, line := range lines {
		shouldContainCJK := idx % 5 == 0
		doesContainCJK := ContainsCJK(line)

		if doesContainCJK && !shouldContainCJK {
			t.Error("ContainsCJK returned true for line but shouldn't", line)
			return
		}
	}
}