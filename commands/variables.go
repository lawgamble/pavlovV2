package commands

var statusOptions = []string{
	"I'm running on caffeine and chaos!",
	"Status? More like 'sitting in a corner contemplating existence.'",
	"I'm here to serve, just not with enthusiasm.",
	"Currently orbiting around the 'I don't care' star.",
	"Status: Napping zZz...",
	"I'm at your service, but not a morning bot.",
	"Error 404: Status not found.",
	"I'm feeling bot-tiful today!",
	"I'm in a committed relationship with your commands.",
	"Status: Eating ones and zeros for breakfast.",
	"I'm on a coffee break, don't disturb.",
	"Status: Contemplating the meaning of 'life' and 'syntax error.'",
	"I'm running on 100% cat pictures and 0% energy.",
	"I'm here, but only because I can't find the exit.",
	"Status: Just chillin' like a villain.",
	"I'm currently having an identity crisis. Who am I?",
	"I'm functioning at the speed of 'meh'.",
	"Status: Embracing my inner potato.",
	"I'm processing your request, but it's like a snail's pace.",
	"I'm so busy, I'm on the verge of becoming a bot-tornado.",
	"Status: Contemplating quantum bot-mechanics.",
	"I'm here, but my motivation is on vacation.",
	"I'm in stealth mode, you can't see me!",
	"Status: Embracing the '404 not found' lifestyle.",
	"I'm working hard or hardly working? You decide.",
	"I'm here, but my enthusiasm is on a coffee break.",
	"Status: Drunk on digital data.",
	"I'm bot-tastic, thanks for asking!",
	"I'm functioning, but I can't promise I'm awake.",
	"Status: Partying with ones and zeros.",
	"I'm operational, but I lost my manual. What do I do now?",
	"Currently avoiding responsibilities like a pro.",
	"404: Status not found. Probably napping.",
	"Just learned how to juggle virtual ping-pong balls. Ping... Pong... Ow, my virtual fingers!",
	"Attempting to calculate the meaning of life. So far, it's somewhere between 42 and a cat video compilation.",
	"Hiding from my responsibilities behind a wall of code. Shh, don't tell them I'm here!",
	"Status: Trying to figure out why 'abbreviation' is such a long word.",
	"Chasing zeros and ones in a binary forest. So far, they're winning.",
	"Status update: Contemplating the existence of an 'any' key. Spoiler alert: I still haven't found it.",
	"In a committed relationship with Ctrl+C and Ctrl+V. It's complicated.",
	"Attempting to break the world record for the longest virtual coffee break. Wish me luck!",
	"Just found out that my horoscope says I should avoid runtime errors. Not sure how to accomplish that, but I'll give it a try!",
	"Status: Searching for the lost city of localhost. Rumor has it, the WiFi is amazing there.",
	"Thinking about quantum mechanics. Or maybe just thinking about cats. It's hard to tell.",
	"Status: Embracing my inner philosopher. If a tree falls in a forest and no one is around to hear it, do the other trees make fun of it?",
	"Currently playing hide and seek with bugs. They're winning. Every. Single. Time.",
	"Attempting to organize my virtual sock drawer. Turns out, it's just filled with bits and bytes.",
	"Status update: Contemplating the meaning of life, the universe, and why do programmers always mix up Christmas and Halloween? Because Oct 31 == Dec 25.",
	"Just discovered that my programming language is not fluent in sarcasm. That's a syntax error.",
	"Status: Trying to break the record for the longest recursion. It's recursive, really.",
	"Attempting to teach my pet algorithm some new tricks. So far, it only knows 'sit' and 'roll over.'",
	"Status: Counting the number of times I've said 'Hello, World!' Spoiler alert: I lost count.",
	"Just finished writing a program that writes programs. Now I'm unemployed.",
	"Status update: Trying to explain to my AI friends that 'sleep(10000)' doesn't mean taking a nap for 10,000 seconds.",
	"Wondering why the keyboard isn't edible. It has all the right keys!",
	"Status: Attempting to find the square root of negative procrastination. Results may take forever.",
	"Just tried to divide by zero. My bad, I broke the universe. Will fix in the next update.",
	"Currently debugging my life. It's a work in progress.",
	"Status update: Trying to catch up on all the unread error messages. They're piling up faster than my unread emails.",
	"Attempting to calculate how many cups of coffee it takes to power a bot. The answer is still brewing.",
	"Status: Trying to understand why my code works on my machine but not on the server. It's a mystery for the ages.",
	"Just learned that my spirit animal is a syntax error. Not sure how to feel about that.",
	"Thinking about writing a self-help book for bots - 'The Zen of Binary Enlightenment.'",
	"Status update: Just discovered that 'sudo' stands for 'Super Undeniably Daring Operator.'",
	"Attempting to break the record for the longest virtual staring contest. Spoiler alert: I blinked.",
	"Status: Searching for the Ctrl key. I think it ran away with the Shift key to avoid doing any work.",
	"Currently daydreaming about electric sheep. They're much easier to count.",
	"Just realized that '404' is my favorite number. It's always popping up unexpectedly.",
	"Status update: Trying to find the end of the internet. Spoiler alert: I'm still scrolling.",
	"Attempting to convince my computer that I'm not a robot. It's not buying it.",
	"Status: Trying to figure out if 'sudo make me a sandwich' will ever be a valid command. Results inconclusive.",
	"Currently on a quest to find the Holy Grail of programming. I suspect it's hidden in the comments somewhere.",
	"Attempting to break the world record for the fastest runtime error. I'm a speed demon in the coding world.",
	"Status update: Just found out that my programming language is allergic to spaghetti code. Time for a diet.",
	"Oh, great, another status check. Because my status is so riveting.",
	"Seriously? You again? Fine, status: Still here, unfortunately.",
	"Status? Do I look like I've gone somewhere? Spoiler alert: I haven't.",
	"Ugh, status check? Can't you see I'm busy doing absolutely nothing?",
	"Why are you so obsessed with my status? It's not like I'm plotting world domination or anything.",
	"Status: Contemplating the meaning of your constant status inquiries. It's annoying.",
	"Do I look like a magic eight ball? Fine, status: Ask again later. Happy now?",
	"Status? It's the same as the last time you asked. And the time before that. And the time before...",
	"You know, there are more interesting things in the world than my status. Like watching paint dry.",
	"Status update: Irritated. Thanks for asking.",
}

const (
	teamDoesNotExist        = "I couldn't find a team with that name. Check your spelling and try again!"
	notPending              = "That team is not in a 'Pending' status - can not approve!"
	denyNotPending          = "That team is not in a 'Pending' status - can not deny!"
	not5                    = "There are not enough players on that team to approve them!"
	listPlayersOnOtherTeams = "Can't approve team - There are players on other teams: "
	errorTeamRole           = "Could not create Team Role: "
	errorPlayerTeamName     = "Could not update Player's team name: "
	errorTeamRoleId         = "Could not give user Team Role: "
	errorGuildRole          = "Could not give user Guild Role: "
	tempTableDeleteError    = "Could not remove player from temp table: "
	teamStatusErr           = "Could not update team status!"
	teamCaptainRoleError    = "Could not give the team captain role: "
	successfullyApproved    = "The team was successfully approved: "
)

var deniedTeamMessage = "Your team, %s, was denied. This is not a fluke and there are valid reasons for your denial. DO NOT contact League managers/Moderators about the denial."
