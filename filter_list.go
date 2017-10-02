package main

import (
	"bytes"
	"strings"
	"unicode"
)

var (
	filterSet = func() map[string]struct{} {
		list := strings.FieldsFunc(filterWordList, func(r rune) bool {
			return r == '\n'
		})
		set := make(map[string]struct{}, len(list))
		empStruct := struct{}{}

		for _, word := range list {
			set[strings.ToLower(word)] = empStruct
		}

		return set
	}()
)

func filterTest(input string) bool {
	words := strings.Fields(input) // TODO: multi word filter ones as well
	for _, word := range words {
		if _, ok := filterSet[word]; ok {
			return true
		}
	}
	return false
}

func filterPrep(input string) string {
	return filterPrepNoLower(strings.ToLower(input))
}

func filterPrepNoLower(input string) string {
	runes := []rune(input)
	buf := bytes.NewBuffer(make([]byte, 0, len(runes)))

	for _, r := range runes {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			buf.WriteRune(r)
		}
	}

	return buf.String()
}

// The big list

const (
	filterWordList = `2 girls 1 cup
2g1c
4r5e
50 yard cunt punt
5h1t
5hit
a2m
a55
a55hole
a_s_s
acrotomophilia
aeolus
ahole
alabama hot pocket
alaskan pipeline
amcik
anal
anal impaler
anal leakage
analprobe
anilingus
anus
apeshit
ar5e
areola
areole
arian
arrse
arse
arsehole
arserape
arsewipe
aryan
ass
ass fuck   
ass hole
ass-fucker
assbang
assbanged
assbangs
asses
assfuck
assfucker
assfukka
assh0le
asshat
assho1e
asshole
assholes
assmaster
assmucus   
assmunch
assramer
assrape
asswhole
asswipe
asswipes
atouche
auto erotic
autoerotic
ayir
azazel
azz
b!tch
b00b
b00bs
b17ch
b1tch
babe
babeland
babes
baby batter
baby juice
bad-ass
badass
ball gag
ball gravy
ball kicking
ball licking
ball sack
ball sucking
ballbag
balls
ballsac
ballsack
bang
bang (one's) box   
bangbros
banger
bareback
barely legal
barenaked
barf
bastard
bastardo
bastards
bastinado
bawdy
bbw
bdsm
beaner
beaners
beardedclam
beastial
beastiality
beastility
beatch
beater
beaver
beaver cleaver
beaver lips
beef curtain   
beer
beeyotch
bellend
benchod
beotch
bestial
bestiality
bi+ch
bi7ch
biatch
big black
big breasts
big knockers
big tits
bigtits
bimbo
bimbos
birdlock
bitch
bitch tit   
bitched
bitcher
bitchers
bitches
bitchin
bitching
bitchy
black cock
blonde action
blonde on blonde action
bloody
blow
blow job
blow me   
blow mud   
blow your load
blowjob
blowjobs
blue waffle
blue waffle   
blumpkin
blumpkin   
bod
bodily
boink
boiolas
bollock
bollocks
bollok
bondage
boned
boner
boners
bong
boob
boobies
boobs
booby
booger
bookie
booobs
boooobs
booooobs
booooooobs
bootee
bootie
booty
booty call
booze
boozer
boozy
bosom
bosomy
bowel
bowels
bra
brassiere
breast
breasts
brown showers
brunette action
buceta
bugger
bukkake
bull shit
bulldyke
bullet vibe
bullshit
bullshits
bullshitted
bullturds
bum
bung
bung hole
bunghole
bunny fucker
busty
butfuck
butt
butt fuck
butt fuck   
buttcheeks
buttfuck
buttfucker
butthole
buttmonkey
buttmuch
buttplug
c-0-c-k
c-o-c-k
c-u-n-t
c.0.c.k
c.o.c.k.
c.u.n.t
c0ck
c0cksucker
cabron
caca
cahone
camel toe
cameltoe
camgirl
camslut
camwhore
carpet muncher
carpetmuncher
cawk
cazzo
cervix
chinc
chincs
chink
choade   
chocolate rosebuds
chode
chodes
chota bags   
chraa
chuj
cipa
circlejerk
cl1t
clamjouster
cleveland steamer
climax
clit
clit licker   
clitoris
clitorus
clits
clitty
clitty litter   
clover clamps
clusterfuck
cnut
cocain
cocaine
cock
cock pocket   
cock snot   
cock sucker
cock-sucker
cockblock
cockface
cockhead
cockholster
cocking
cockknocker
cockmunch
cockmuncher
cocks
cockslap
cockslapped
cockslapping
cocksmoker
cocksuck
cocksucked
cocksucker
cocksucking
cocksucks
cocksuka
cocksukka
coital
cok
cokmuncher
coksucka
condom
coon
coons
coprolagnia
coprophilia
corksucker
cornhole
cornhole   
corp whore   
cox
crabs
crackwhore
creampie
cum
cum chugger   
cum dumpster   
cum freak   
cum guzzler   
cumdump   
cummer
cummin
cumming
cums
cumshot
cumshots
cumslut
cumstain
cunilingus
cunillingus
cunnilingus
cunny
cunt
cunt hair   
cunt-struck   
cuntalot
cuntbag   
cuntface
cuntfish
cunthunter
cunting
cuntlick
cuntlicker
cuntlicking
cuntree
cunts
cuntsicle   
cut rope   
cyalis
cyberfuc
cyberfuck
cyberfucked
cyberfucker
cyberfuckers
cyberfucking
d0ng
d0uch3
d0uche
d1ck
d1ld0
d1ldo
d4mn
dago
dagos
dammit
damn
damned
damnit
damnnation
darkie
date rape
daterape
dawgie-style
daygo
dayum
deep throat
deepthroat
dego
dendrophilia
dick
dick hole   
dick shy   
dick-ish
dickabout
dickaround
dickbag
dickdipper
dickface
dickflipper
dickhead
dickheads
dicking
dickish
dickripper
dicksipper
dickwad
dickward
dickweed
dickwhipper
dickzipper
diddle
dike
dildo
dildos
diligaf
dillweed
dimwit
dingle
dingleberries
dingleberry
dink
dinks
dipship
dirsa
dirty pillows
dirty sanchez
dirty sanchez   
dlck
dog style
dog-fucker
doggie style
doggie-style
doggiestyle
doggin
dogging
doggy style
doggy-style
doggystyle
dolcett
domination
dominatrix
dommes
dong
donkey punch
donkeyribber
doosh
dopey
double dong
double penetration
douch3
douche
douchebag
douchebags
douchey
dp action
drunk
dry hump
duche
dumass
dumbass
dumbasses
dupa
dvda
dyke
dykes
dziwka
eat a dick   
eat hair pie   
eat my ass
ecchi
ejaculate
ejaculated
ejaculates
ejaculating
ejaculatings
ejaculation
ejakulate
ekrem
ekto
enculer
enlargement
erect
erection
erotic
erotism
escort
essohbee
eunuch
extacy
extasy
f u c k
f u c k e r
f-u-c-k
f.u.c.k
f4nny
f_u_c_k
facial   
fack
faen
fag
fagg
fagged
fagging
faggit
faggitt
faggot
faggs
fagot
fagots
fags
faig
faigt
fancul
fanny
fannybandit
fannyflaps
fannyfucker
fanyy
fart
farted
farting
fartings
fartknocker
farts
farty
fatass
fcuk
fcuker
fcuking
fecal
feces
feck
fecker
felatio
felch
felcher
felching
fellate
fellatio
feltch
feltcher
female squirting
femdom
ficken
figging
fingerbang
fingerfuck
fingerfucked
fingerfucker
fingerfuckers
fingerfucking
fingerfucks
fingering
fist fuck   
fisted
fistfuck
fistfucked
fistfucker
fistfuckers
fistfucking
fistfuckings
fistfucks
fisting
fisty
fitta
fitte
flange
flikker
flog the log   
floozy
foad
fondle
foobar
fook
fooker
foot fetish
footjob
foreskin
fotze
freex
frigg
frigga
frotting
ftq
fubar
fuck
fuck buttons
fuck hole   
fuck puppet   
fuck trophy   
fuck yo mama   
fuck-ass   
fuck-bitch   
fuck-tard
fucka
fuckass
fucked
fucker
fuckers
fuckface
fuckhead
fuckheads
fuckin
fucking
fuckings
fuckingshitmotherfucker
fuckmaster
fuckme
fuckmeat   
fucknugget
fucknut
fuckoff
fucks
fucktard
fucktards
fucktoy   
fuckup
fuckwad
fuckwhit
fuckwit
fucky
fuck   
fudge packer
fudgepacker
fuk
fuker
fukker
fukkin
fuks
fukwhit
fukwit
futanari
futkretzn
fux
fux0r
fvck
fxck
g-spot
gae
gai
gang bang
gang-bang   
gangbang
gangbanged
gangbangs
gangbang   
ganja
gash
gassy ass   
gay
gay sex
gaylord
gays
gaysex
genitals
gey
gfy
ghay
ghey
giant cock
gigolo
girl on
girl on top
girls gone wild
glans
goatcx
goatse
god
god damn
god-dam
god-damned
godamn
godamnit
goddam
goddammit
goddamn
goddamned
gokkun
golden shower
goldenshower
gonad
gonads
goo girl
goodpoop
gook
gooks
goolies
goregasm
gringo
grope
group sex
gspot
gtfo
guido
guiena
guro
h0m0
h0mo
h0r
h4x0r
ham flap   
hand job
handjob
hard core
hard on
hardcore
hardcoresex
he11
hebe
heeb
hell
helvete
hentai
heroin
herp
herpes
herpy
heshe
hitler
hiv
hoar
hoare
hobag
hoer
hom0
homey
homo
homoerotic
homoey
honkey
honky
hooch
hookah
hooker
hoor
hootch
hooter
hooters
hore
horniest
horny
hot carl
hot chick
hotsex
how to kill
how to murdep
how to murder
huevon
huge fat
hui
hump
humped
humping
hussy
hymen
inbred
incest
injun
intercourse
j3rk0ff
jack off
jack-off
jackass
jackhole
jackoff
jail bait
jailbait
jap
japs
jelly donut
jerk
jerk off
jerk-off
jerk0ff
jerked
jerkoff
jigaboo
jiggaboo
jiggerboo
jism
jiz
jizm
jizz
jizzed
juggs
kaffir
kawk
kike
kikes
kinbaku
kinkster
kinky
kinky jesus   
kkk
klan
knob
knob end
knobbing
knobead
knobed
knobend
knobhead
knobjockey
knobjocky
knobjokey
knulle
kock
kondum
kondums
kooch
kooches
kootch
kraut
kuk
kuksuger
kum
kumer
kummer
kumming
kums
kunilingus
kurac
kurwa
kusi
kwif   
kyke
kyrp
l3i+ch
l3itch
labia
leather restraint
leather straight jacket
lech
lemon party
len
leper
lesbian
lesbians
lesbo
lesbos
lez
lezbian
lezbians
lezbo
lezbos
lezzie
lezzies
lezzy
lmao
lmfao
loin
loins
lolita
lovemaking
lube
lust
lusting
lusty
m-fucking
m0f0
m0fo
m45terbate
ma5terb8
ma5terbate
mafugly   
make me come
male squirting
mamhoon
mams
masochist
massa
master-bate
masterb8
masterbat*
masterbat3
masterbate
masterbating
masterbation
masterbations
masturbat
masturbate
masturbating
masturbation
maxi
menage a trois
menses
menstruate
menstruation
merd
merde
meth
mibun
milf
minge
minger
missionary position
mo-fo
mof0
mofo
molest
mong
monkleigh
moolie
mothafuck
mothafucka
mothafuckas
mothafuckaz
mothafucked
mothafucker
mothafuckers
mothafuckin
mothafucking
mothafuckings
mothafucks
mother fucker
mother fucker   
motherfuck
motherfucka
motherfucked
motherfucker
motherfuckers
motherfuckin
motherfucking
motherfuckings
motherfuckka
motherfucks
mouliewop
mound of venus
mr hands
mtherfucker
mthrfucker
mthrfucking
muff
muff diver
muff puff   
muffdiver
muffdiving
muffmuncher
muie
mulkku
mummyporn
munter
murder
muschi
mutha
muthafecker
muthafuckaz
muthafucker
muthafuckker
muther
mutherfucker
mutherfucking
muthrfucking
n1gga
n1gger
nad
nads
naked
nambla
napalm
nappy
nawashi
nazi
nazis
nazism
need the dick   
negro
neonazi
nepesaurio
nig nog
niger
nigg3r
nigg4h
nigga
niggah
niggar
niggars
niggas
niggaz
nigger
niggers
niggle
niglet
nimphomania
nimrod
ninny
nipple
nipples
nob
nob jokey
nobhead
nobjocky
nobjokey
nooky
nsfw
nsfw images
nude
nudes
nudity
numbnuts
nut butter   
nutsack
nympho
nymphomania
octopussy
omorashi
one cup two girls
one guy one jar
ootzak
opiate
opium
organ
orgasim
orgasims
orgasm
orgasmic
orgasms
orgies
orgy
orospu
ovary
ovum
ovums
p.u.s.s.y
p.u.s.s.y.
p0rn
p0rnhub
paddy
paedophile
paki
panties
panty
paska
pastie
pasty
pawn
pcp
pecker
pedo
pedobear
pedophile
pedophilia
pedophiliac
pee
peepee
pegging
pendejo
penetrate
penetration
penial
penile
penis
penisfucker
penisperse
perversion
peyote
phalli
phallic
phone sex
phonesex
phuck
phuk
phuked
phuking
phukked
phukking
phuks
phuq
picka
piece of shit
pierdol
pigfucker
pikey
pillowbiter
pillu
pimmel
pimp
pimpis
pinko
pis
pises
pisin
pising
pisof
piss
piss pig
piss-off
pissed
pisser
pissers
pisses
pissflaps
pissin
pissing
pissoff
pisspig
pizdapoontsee
playboy
pleasure chest
pms
polack
pole smoker
pollock
ponyplay
poof
poon
poontang
poop
poop chute
poopchute
porn
pornhub
porno
pornography
pornos
pos
pot
potty
pr0n
pr0nhub
preteen
preud
prick
pricks
prig
prince albert piercing
pron
pronhub
prostitute
prude
pthc
pube
pubes
pubic
pubis
pula
pule
punany
punkass
punky
pusies
puss
pusse
pussi
pussies
pussy
pussy fart   
pussy palace   
pussypounder
pussys
pusy
pusys
puta
puto
qahbeh
queaf
queaf   
queef
queer
queero
queers
quicky
quim
qweef
r-tard
r3dtub3
r3dtube
racy
raghead
raging boner
rape
raped
raper
raping
rapist
raunch
rautenbergschaffer
rectal
rectum
rectus
redtub3
redtube
reefer
reetard
reich
retard
retarded
reverse cowgirl
revue
rimjaw
rimjob
rimming
ritard
rosy palm
rosy palm and her 5 sisters
rtard
rum
rump
rumprammer
ruski
rusty trombone
s hit
s&m
s-h-1-t
s-h-i-t
s-o-b
s.h.i.t.
s.o.b.
s0b
s_h_i_t
sadism
sadist
sandbar   
santorum
sausage queen   
scag
scantily
scat
scheiss
scheisse
schizo
schlampe
schlong
schmuck
scissoring
screw
screwed
screwing
scroat
scrog
scrot
scrote
scrotum
scrud
scum
seaman
seamen
seduce
semen
sex
sexeh
sexist
sexo
sexual
sexy
sh!+
sh!t
sh1t
shag
shagged
shagger
shaggin
shagging
shamedame
sharmuta
sharmute
shaved beaver
shaved pussy
shemale
shenzi
shi+
shiat
shibari
shipal
shit
shit fucker   
shitblimp
shitdick
shite
shiteater
shited
shitey
shitface
shitfuck
shitfull
shithead
shithole
shithouse
shiting
shitings
shits
shitt
shitted
shitter
shitters
shitting
shittings
shitty
shity
shiz
shizer
shota
shrimping
skag
skank
skeet
skribz
skurwysyn
slag
slanteye
slave
sleaze
sleazy
slope   
slut
slut bucket   
slutdumper
slutkiss
sluts
smegma
smut
smutty
snatch
sniper
snowballing
snuff
sodding
sodom
sodomize
sodomy
son-of-a-bitch
souse
soused
spac
spacker
spacko
spank
spastic
spaz
sperm
sphencter
spic
spick
spierdalaj
spik
spiks
splooge
splooge moose
spooge
spread legs
spunk
spunking
steamy
stfu
stiffy
stoned
strap on
strapon
strappado
strip
strip club
stroke
style doggy
sucking
suicide girls
suka
sultry women
sumofabiatch
swastika
swinger
t1t
t1tt1e5
t1tties
tainted love
tampon
tard
taste my
tawdry
tea bagging
teabagging
teat
teets
teez
terd
teste
testee
testes
testical
testicle
testis
threesome
throating
thrust
thug
tied up
tight white
tinkle
tit
tit wank   
titfuck
titi
tits
titt
tittie5
tittiefucker
titties
titty
tittyfuck
tittyfucker
tittywank
titwank
toke
tongue in a
toots
topless
tosser
towelhead
tramp
tranny
transsexual
trashy
tribadism
tub girl
tubgirl
turd
tush
tushy
tw4t
twat
twathead
twats
twatty
twink
twinkie
two girls one cup
twunt
twunter
ugly
uncunt
undies
undressing
unwed
upskirt
urethra play
urinal
urine
urophilia
uterus
uzi
v14gra
v1gra
vag
vagina
valium
venus mound
viagra
violet wand
virgin
vixen
vodka
vorarephilia
voyeur
vulgar
vulva
w00se
wad
wang
wank
wanked
wanker
wankered
wankers
wanking
wanky
wazoo
wedgie
weed
weenie
weewee
weiner
weirdo
wench
wet dream
wetback
wh0re
wh0reface
white power
whitey
whiz
whoar
whoralicious
whore
whorealicious
whored
whoreface
whorehopper
whorehouse
whores
whoring
wigger
willies
willy
wog
womb
woody
wop
wrapping men
wrinkled starfish
wtf
x-rated
xrated
xxx
yaoi
yeasty
yellow showers
yiffy
yobbo
zoophile
zoophilia
french kisses
french kiss
fapping
breaking up`
)
