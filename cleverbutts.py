#!/usr/bin/env python3
import asyncio
import async_timeout
import signal
import random
import json
import discord

CHATENGINE_URL = 'https://chatengine.xyz/ask'

f = open('keys.json', 'rb')
keydat = json.loads(f.read().decode('utf-8'))
CHATENGINE_KEY = keydat['api'][1]

tokens = keydat['discord']
active = True

try:
    import uvloop
    asyncio.set_event_loop_policy(uvloop.EventLoopPolicy())
except ImportError:
    pass
loop = asyncio.get_event_loop()

class Cleverbutt(discord.Client):
    def __init__(self):
        discord.Client.__init__(self)

        self.cleverbutt_timers = set()
        self.cleverbutt_latest = {}
        self.cleverbutt_replied_to = set()

    async def on_ready(self):
        print('{} ready'.format(self.user))
        self.clever_chan = self.get_channel(333718439925383189)

    async def askcb(self, query: str) -> str:
        with async_timeout.timeout(5):
            async with self.http._session.post(CHATENGINE_URL, headers={'Referer': CHATENGINE_KEY}, data=query) as resp:
                return (await resp.read()).decode('utf-8')

    async def clever_reply(self, msg: discord.Message):
        """Cleverbutts handler."""
        self.cleverbutt_timers.add(msg.guild.id)
        await asyncio.sleep(random.random())
        with msg.channel.typing():
            await asyncio.sleep(random.random() * 1.8)
            try:
                query = self.cleverbutt_latest[msg.guild.id]
            except KeyError:
                query = msg.content
            reply_bot = await self.askcb(query)
            s_duration = (((len(reply_bot) / 15) * 1.4) + random.random()) - 0.2
            await asyncio.sleep(s_duration / 1.5)
            await msg.channel.send(reply_bot)
        await asyncio.sleep(0.5)
        try:
            del self.cleverbutt_latest[msg.guild.id]
        except Exception:
            pass
        self.cleverbutt_replied_to.add(msg.id)
        self.cleverbutt_timers.remove(msg.guild.id)

    async def on_pm(self, msg: discord.Message):
        """PM replying logic."""
        if msg.content.startswith('`'):
            if 'REPL' in self.bot.cogs:
                if msg.channel.id in self.bot.cogs['REPL'].sessions:
                    return
        await msg.channel.trigger_typing()
        c = msg.content
        for m in msg.mentions:
            c = c.replace(m.mention, m.display_name)
        cb_reply = await self.askcb(c)
        return await msg.channel.send(':speech_balloon: ' + cb_reply)

    async def on_message(self, message: discord.Message):
        if message.author == self.user: return
        if isinstance(message.channel, discord.abc.PrivateChannel):
            await self.on_pm(msg)
            return
        if not message.author.bot:
            if message.author.id == 160567046642335746 and message.content.lower().startswith('c::') and self.user.id == 333715715926523914:
                cmd = message.content.lower()[3:]

                if cmd == 'help':
                    await message.channel.send('Commands: `kickstart`, `toggle`, `exit`')
                elif cmd == 'kickstart':
                    await self.clever_chan.send('Hello there!')
                    await message.channel.send('Done.')
                elif cmd == 'toggle':
                    active = not active
                    await message.channel.send('active: ' + str(active))
                elif cmd == 'exit':
                    await message.channel.send('Bye')
                    loop.stop()
            return
        if message.channel != self.clever_chan: return

        if message.guild.id in self.cleverbutt_timers: # still on timer for next response
            self.cleverbutt_latest[message.guild.id] = message.content
        else:
            await self.clever_reply(message)

bots = []
for token in tokens:
    bot = Cleverbutt()
    loop.create_task(bot.start(token))
    bots.append(bot)

signal.signal(signal.SIGINT, lambda: loop.stop())
loop.run_forever()
