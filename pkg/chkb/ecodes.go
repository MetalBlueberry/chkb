/*
Copyright © 2020 Víctor Pérez @MetalBlueberry

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package chkb

const (
	EV_SYN       = 0x00
	EV_KEY       = 0x01
	EV_REL       = 0x02
	EV_ABS       = 0x03
	EV_MSC       = 0x04
	EV_SW        = 0x05
	EV_LED       = 0x11
	EV_SND       = 0x12
	EV_REP       = 0x14
	EV_FF        = 0x15
	EV_PWR       = 0x16
	EV_FF_STATUS = 0x17
	EV_MAX       = 0x1f
)

const (
	KEY_RESERVED                 KeyCode = 0
	KEY_ESC                      KeyCode = 1
	KEY_1                        KeyCode = 2
	KEY_2                        KeyCode = 3
	KEY_3                        KeyCode = 4
	KEY_4                        KeyCode = 5
	KEY_5                        KeyCode = 6
	KEY_6                        KeyCode = 7
	KEY_7                        KeyCode = 8
	KEY_8                        KeyCode = 9
	KEY_9                        KeyCode = 10
	KEY_0                        KeyCode = 11
	KEY_MINUS                    KeyCode = 12
	KEY_EQUAL                    KeyCode = 13
	KEY_BACKSPACE                KeyCode = 14
	KEY_TAB                      KeyCode = 15
	KEY_Q                        KeyCode = 16
	KEY_W                        KeyCode = 17
	KEY_E                        KeyCode = 18
	KEY_R                        KeyCode = 19
	KEY_T                        KeyCode = 20
	KEY_Y                        KeyCode = 21
	KEY_U                        KeyCode = 22
	KEY_I                        KeyCode = 23
	KEY_O                        KeyCode = 24
	KEY_P                        KeyCode = 25
	KEY_LEFTBRACE                KeyCode = 26
	KEY_RIGHTBRACE               KeyCode = 27
	KEY_ENTER                    KeyCode = 28
	KEY_LEFTCTRL                 KeyCode = 29
	KEY_A                        KeyCode = 30
	KEY_S                        KeyCode = 31
	KEY_D                        KeyCode = 32
	KEY_F                        KeyCode = 33
	KEY_G                        KeyCode = 34
	KEY_H                        KeyCode = 35
	KEY_J                        KeyCode = 36
	KEY_K                        KeyCode = 37
	KEY_L                        KeyCode = 38
	KEY_SEMICOLON                KeyCode = 39
	KEY_APOSTROPHE               KeyCode = 40
	KEY_GRAVE                    KeyCode = 41
	KEY_LEFTSHIFT                KeyCode = 42
	KEY_BACKSLASH                KeyCode = 43
	KEY_Z                        KeyCode = 44
	KEY_X                        KeyCode = 45
	KEY_C                        KeyCode = 46
	KEY_V                        KeyCode = 47
	KEY_B                        KeyCode = 48
	KEY_N                        KeyCode = 49
	KEY_M                        KeyCode = 50
	KEY_COMMA                    KeyCode = 51
	KEY_DOT                      KeyCode = 52
	KEY_SLASH                    KeyCode = 53
	KEY_RIGHTSHIFT               KeyCode = 54
	KEY_KPASTERISK               KeyCode = 55
	KEY_LEFTALT                  KeyCode = 56
	KEY_SPACE                    KeyCode = 57
	KEY_CAPSLOCK                 KeyCode = 58
	KEY_F1                       KeyCode = 59
	KEY_F2                       KeyCode = 60
	KEY_F3                       KeyCode = 61
	KEY_F4                       KeyCode = 62
	KEY_F5                       KeyCode = 63
	KEY_F6                       KeyCode = 64
	KEY_F7                       KeyCode = 65
	KEY_F8                       KeyCode = 66
	KEY_F9                       KeyCode = 67
	KEY_F10                      KeyCode = 68
	KEY_NUMLOCK                  KeyCode = 69
	KEY_SCROLLLOCK               KeyCode = 70
	KEY_KP7                      KeyCode = 71
	KEY_KP8                      KeyCode = 72
	KEY_KP9                      KeyCode = 73
	KEY_KPMINUS                  KeyCode = 74
	KEY_KP4                      KeyCode = 75
	KEY_KP5                      KeyCode = 76
	KEY_KP6                      KeyCode = 77
	KEY_KPPLUS                   KeyCode = 78
	KEY_KP1                      KeyCode = 79
	KEY_KP2                      KeyCode = 80
	KEY_KP3                      KeyCode = 81
	KEY_KP0                      KeyCode = 82
	KEY_KPDOT                    KeyCode = 83
	KEY_ZENKAKUHANKAKU           KeyCode = 85
	KEY_102ND                    KeyCode = 86
	KEY_F11                      KeyCode = 87
	KEY_F12                      KeyCode = 88
	KEY_RO                       KeyCode = 89
	KEY_KATAKANA                 KeyCode = 90
	KEY_HIRAGANA                 KeyCode = 91
	KEY_HENKAN                   KeyCode = 92
	KEY_KATAKANAHIRAGANA         KeyCode = 93
	KEY_MUHENKAN                 KeyCode = 94
	KEY_KPJPCOMMA                KeyCode = 95
	KEY_KPENTER                  KeyCode = 96
	KEY_RIGHTCTRL                KeyCode = 97
	KEY_KPSLASH                  KeyCode = 98
	KEY_SYSRQ                    KeyCode = 99
	KEY_RIGHTALT                 KeyCode = 100
	KEY_LINEFEED                 KeyCode = 101
	KEY_HOME                     KeyCode = 102
	KEY_UP                       KeyCode = 103
	KEY_PAGEUP                   KeyCode = 104
	KEY_LEFT                     KeyCode = 105
	KEY_RIGHT                    KeyCode = 106
	KEY_END                      KeyCode = 107
	KEY_DOWN                     KeyCode = 108
	KEY_PAGEDOWN                 KeyCode = 109
	KEY_INSERT                   KeyCode = 110
	KEY_DELETE                   KeyCode = 111
	KEY_MACRO                    KeyCode = 112
	KEY_MUTE                     KeyCode = 113
	KEY_VOLUMEDOWN               KeyCode = 114
	KEY_VOLUMEUP                 KeyCode = 115
	KEY_POWER                    KeyCode = 116
	KEY_KPEQUAL                  KeyCode = 117
	KEY_KPPLUSMINUS              KeyCode = 118
	KEY_PAUSE                    KeyCode = 119
	KEY_SCALE                    KeyCode = 120
	KEY_KPCOMMA                  KeyCode = 121
	KEY_HANGEUL                  KeyCode = 122
	KEY_HANGUEL                  KeyCode = 123
	KEY_HANJA                    KeyCode = 123
	KEY_YEN                      KeyCode = 124
	KEY_LEFTMETA                 KeyCode = 125
	KEY_RIGHTMETA                KeyCode = 126
	KEY_COMPOSE                  KeyCode = 127
	KEY_STOP                     KeyCode = 128
	KEY_AGAIN                    KeyCode = 129
	KEY_PROPS                    KeyCode = 130
	KEY_UNDO                     KeyCode = 131
	KEY_FRONT                    KeyCode = 132
	KEY_COPY                     KeyCode = 133
	KEY_OPEN                     KeyCode = 134
	KEY_PASTE                    KeyCode = 135
	KEY_FIND                     KeyCode = 136
	KEY_CUT                      KeyCode = 137
	KEY_HELP                     KeyCode = 138
	KEY_MENU                     KeyCode = 139
	KEY_CALC                     KeyCode = 140
	KEY_SETUP                    KeyCode = 141
	KEY_SLEEP                    KeyCode = 142
	KEY_WAKEUP                   KeyCode = 143
	KEY_FILE                     KeyCode = 144
	KEY_SENDFILE                 KeyCode = 145
	KEY_DELETEFILE               KeyCode = 146
	KEY_XFER                     KeyCode = 147
	KEY_PROG1                    KeyCode = 148
	KEY_PROG2                    KeyCode = 149
	KEY_WWW                      KeyCode = 150
	KEY_MSDOS                    KeyCode = 151
	KEY_COFFEE                   KeyCode = 152
	KEY_SCREENLOCK               KeyCode = 153
	KEY_ROTATE_DISPLAY           KeyCode = 153
	KEY_DIRECTION                KeyCode = 154
	KEY_CYCLEWINDOWS             KeyCode = 154
	KEY_MAIL                     KeyCode = 155
	KEY_BOOKMARKS                KeyCode = 156
	KEY_COMPUTER                 KeyCode = 157
	KEY_BACK                     KeyCode = 158
	KEY_FORWARD                  KeyCode = 159
	KEY_CLOSECD                  KeyCode = 160
	KEY_EJECTCD                  KeyCode = 161
	KEY_EJECTCLOSECD             KeyCode = 162
	KEY_NEXTSONG                 KeyCode = 163
	KEY_PLAYPAUSE                KeyCode = 164
	KEY_PREVIOUSSONG             KeyCode = 165
	KEY_STOPCD                   KeyCode = 166
	KEY_RECORD                   KeyCode = 167
	KEY_REWIND                   KeyCode = 168
	KEY_PHONE                    KeyCode = 169
	KEY_ISO                      KeyCode = 170
	KEY_CONFIG                   KeyCode = 171
	KEY_HOMEPAGE                 KeyCode = 172
	KEY_REFRESH                  KeyCode = 173
	KEY_EXIT                     KeyCode = 174
	KEY_MOVE                     KeyCode = 175
	KEY_EDIT                     KeyCode = 176
	KEY_SCROLLUP                 KeyCode = 177
	KEY_SCROLLDOWN               KeyCode = 178
	KEY_KPLEFTPAREN              KeyCode = 179
	KEY_KPRIGHTPAREN             KeyCode = 180
	KEY_NEW                      KeyCode = 181
	KEY_REDO                     KeyCode = 182
	KEY_F13                      KeyCode = 183
	KEY_F14                      KeyCode = 184
	KEY_F15                      KeyCode = 185
	KEY_F16                      KeyCode = 186
	KEY_F17                      KeyCode = 187
	KEY_F18                      KeyCode = 188
	KEY_F19                      KeyCode = 189
	KEY_F20                      KeyCode = 190
	KEY_F21                      KeyCode = 191
	KEY_F22                      KeyCode = 192
	KEY_F23                      KeyCode = 193
	KEY_F24                      KeyCode = 194
	KEY_PLAYCD                   KeyCode = 200
	KEY_PAUSECD                  KeyCode = 201
	KEY_PROG3                    KeyCode = 202
	KEY_PROG4                    KeyCode = 203
	KEY_DASHBOARD                KeyCode = 204
	KEY_SUSPEND                  KeyCode = 205
	KEY_CLOSE                    KeyCode = 206
	KEY_PLAY                     KeyCode = 207
	KEY_FASTFORWARD              KeyCode = 208
	KEY_BASSBOOST                KeyCode = 209
	KEY_PRINT                    KeyCode = 210
	KEY_HP                       KeyCode = 211
	KEY_CAMERA                   KeyCode = 212
	KEY_SOUND                    KeyCode = 213
	KEY_QUESTION                 KeyCode = 214
	KEY_EMAIL                    KeyCode = 215
	KEY_CHAT                     KeyCode = 216
	KEY_SEARCH                   KeyCode = 217
	KEY_CONNECT                  KeyCode = 218
	KEY_FINANCE                  KeyCode = 219
	KEY_SPORT                    KeyCode = 220
	KEY_SHOP                     KeyCode = 221
	KEY_ALTERASE                 KeyCode = 222
	KEY_CANCEL                   KeyCode = 223
	KEY_BRIGHTNESSDOWN           KeyCode = 224
	KEY_BRIGHTNESSUP             KeyCode = 225
	KEY_MEDIA                    KeyCode = 226
	KEY_SWITCHVIDEOMODE          KeyCode = 227
	KEY_KBDILLUMTOGGLE           KeyCode = 228
	KEY_KBDILLUMDOWN             KeyCode = 229
	KEY_KBDILLUMUP               KeyCode = 230
	KEY_SEND                     KeyCode = 231
	KEY_REPLY                    KeyCode = 232
	KEY_FORWARDMAIL              KeyCode = 233
	KEY_SAVE                     KeyCode = 234
	KEY_DOCUMENTS                KeyCode = 235
	KEY_BATTERY                  KeyCode = 236
	KEY_BLUETOOTH                KeyCode = 237
	KEY_WLAN                     KeyCode = 238
	KEY_UWB                      KeyCode = 239
	KEY_UNKNOWN                  KeyCode = 240
	KEY_VIDEO_NEXT               KeyCode = 241
	KEY_VIDEO_PREV               KeyCode = 242
	KEY_BRIGHTNESS_CYCLE         KeyCode = 243
	KEY_BRIGHTNESS_AUTO          KeyCode = 244
	KEY_DISPLAY_OFF              KeyCode = 245
	KEY_WWAN                     KeyCode = 246
	KEY_RFKILL                   KeyCode = 247
	KEY_MICMUTE                  KeyCode = 248
	BTN_MISC                     KeyCode = 0x100
	BTN_0                        KeyCode = 0x100
	BTN_1                        KeyCode = 0x101
	BTN_2                        KeyCode = 0x102
	BTN_3                        KeyCode = 0x103
	BTN_4                        KeyCode = 0x104
	BTN_5                        KeyCode = 0x105
	BTN_6                        KeyCode = 0x106
	BTN_7                        KeyCode = 0x107
	BTN_8                        KeyCode = 0x108
	BTN_9                        KeyCode = 0x109
	BTN_MOUSE                    KeyCode = 0x110
	BTN_LEFT                     KeyCode = 0x110
	BTN_RIGHT                    KeyCode = 0x111
	BTN_MIDDLE                   KeyCode = 0x112
	BTN_SIDE                     KeyCode = 0x113
	BTN_EXTRA                    KeyCode = 0x114
	BTN_FORWARD                  KeyCode = 0x115
	BTN_BACK                     KeyCode = 0x116
	BTN_TASK                     KeyCode = 0x117
	BTN_JOYSTICK                 KeyCode = 0x120
	BTN_TRIGGER                  KeyCode = 0x120
	BTN_THUMB                    KeyCode = 0x121
	BTN_THUMB2                   KeyCode = 0x122
	BTN_TOP                      KeyCode = 0x123
	BTN_TOP2                     KeyCode = 0x124
	BTN_PINKIE                   KeyCode = 0x125
	BTN_BASE                     KeyCode = 0x126
	BTN_BASE2                    KeyCode = 0x127
	BTN_BASE3                    KeyCode = 0x128
	BTN_BASE4                    KeyCode = 0x129
	BTN_BASE5                    KeyCode = 0x12a
	BTN_BASE6                    KeyCode = 0x12b
	BTN_DEAD                     KeyCode = 0x12f
	BTN_GAMEPAD                  KeyCode = 0x130
	BTN_SOUTH                    KeyCode = 0x130
	BTN_A                        KeyCode = BTN_SOUTH
	BTN_EAST                     KeyCode = 0x131
	BTN_B                        KeyCode = BTN_EAST
	BTN_C                        KeyCode = 0x132
	BTN_NORTH                    KeyCode = 0x133
	BTN_X                        KeyCode = BTN_NORTH
	BTN_WEST                     KeyCode = 0x134
	BTN_Y                        KeyCode = BTN_WEST
	BTN_Z                        KeyCode = 0x135
	BTN_TL                       KeyCode = 0x136
	BTN_TR                       KeyCode = 0x137
	BTN_TL2                      KeyCode = 0x138
	BTN_TR2                      KeyCode = 0x139
	BTN_SELECT                   KeyCode = 0x13a
	BTN_START                    KeyCode = 0x13b
	BTN_MODE                     KeyCode = 0x13c
	BTN_THUMBL                   KeyCode = 0x13d
	BTN_THUMBR                   KeyCode = 0x13e
	BTN_DIGI                     KeyCode = 0x140
	BTN_TOOL_PEN                 KeyCode = 0x140
	BTN_TOOL_RUBBER              KeyCode = 0x141
	BTN_TOOL_BRUSH               KeyCode = 0x142
	BTN_TOOL_PENCIL              KeyCode = 0x143
	BTN_TOOL_AIRBRUSH            KeyCode = 0x144
	BTN_TOOL_FINGER              KeyCode = 0x145
	BTN_TOOL_MOUSE               KeyCode = 0x146
	BTN_TOOL_LENS                KeyCode = 0x147
	BTN_TOOL_QUINTTAP            KeyCode = 0x148
	BTN_TOUCH                    KeyCode = 0x14a
	BTN_STYLUS                   KeyCode = 0x14b
	BTN_STYLUS2                  KeyCode = 0x14c
	BTN_TOOL_DOUBLETAP           KeyCode = 0x14d
	BTN_TOOL_TRIPLETAP           KeyCode = 0x14e
	BTN_TOOL_QUADTAP             KeyCode = 0x14f
	BTN_WHEEL                    KeyCode = 0x150
	BTN_GEAR_DOWN                KeyCode = 0x150
	BTN_GEAR_UP                  KeyCode = 0x151
	KEY_OK                       KeyCode = 0x160
	KEY_SELECT                   KeyCode = 0x161
	KEY_GOTO                     KeyCode = 0x162
	KEY_CLEAR                    KeyCode = 0x163
	KEY_POWER2                   KeyCode = 0x164
	KEY_OPTION                   KeyCode = 0x165
	KEY_INFO                     KeyCode = 0x166
	KEY_TIME                     KeyCode = 0x167
	KEY_VENDOR                   KeyCode = 0x168
	KEY_ARCHIVE                  KeyCode = 0x169
	KEY_PROGRAM                  KeyCode = 0x16a
	KEY_CHANNEL                  KeyCode = 0x16b
	KEY_FAVORITES                KeyCode = 0x16c
	KEY_EPG                      KeyCode = 0x16d
	KEY_PVR                      KeyCode = 0x16e
	KEY_MHP                      KeyCode = 0x16f
	KEY_LANGUAGE                 KeyCode = 0x170
	KEY_TITLE                    KeyCode = 0x171
	KEY_SUBTITLE                 KeyCode = 0x172
	KEY_ANGLE                    KeyCode = 0x173
	KEY_ZOOM                     KeyCode = 0x174
	KEY_MODE                     KeyCode = 0x175
	KEY_KEYBOARD                 KeyCode = 0x176
	KEY_SCREEN                   KeyCode = 0x177
	KEY_PC                       KeyCode = 0x178
	KEY_TV                       KeyCode = 0x179
	KEY_TV2                      KeyCode = 0x17a
	KEY_VCR                      KeyCode = 0x17b
	KEY_VCR2                     KeyCode = 0x17c
	KEY_SAT                      KeyCode = 0x17d
	KEY_SAT2                     KeyCode = 0x17e
	KEY_CD                       KeyCode = 0x17f
	KEY_TAPE                     KeyCode = 0x180
	KEY_RADIO                    KeyCode = 0x181
	KEY_TUNER                    KeyCode = 0x182
	KEY_PLAYER                   KeyCode = 0x183
	KEY_TEXT                     KeyCode = 0x184
	KEY_DVD                      KeyCode = 0x185
	KEY_AUX                      KeyCode = 0x186
	KEY_MP3                      KeyCode = 0x187
	KEY_AUDIO                    KeyCode = 0x188
	KEY_VIDEO                    KeyCode = 0x189
	KEY_DIRECTORY                KeyCode = 0x18a
	KEY_LIST                     KeyCode = 0x18b
	KEY_MEMO                     KeyCode = 0x18c
	KEY_CALENDAR                 KeyCode = 0x18d
	KEY_RED                      KeyCode = 0x18e
	KEY_GREEN                    KeyCode = 0x18f
	KEY_YELLOW                   KeyCode = 0x190
	KEY_BLUE                     KeyCode = 0x191
	KEY_CHANNELUP                KeyCode = 0x192
	KEY_CHANNELDOWN              KeyCode = 0x193
	KEY_FIRST                    KeyCode = 0x194
	KEY_LAST                     KeyCode = 0x195
	KEY_AB                       KeyCode = 0x196
	KEY_NEXT                     KeyCode = 0x197
	KEY_RESTART                  KeyCode = 0x198
	KEY_SLOW                     KeyCode = 0x199
	KEY_SHUFFLE                  KeyCode = 0x19a
	KEY_BREAK                    KeyCode = 0x19b
	KEY_PREVIOUS                 KeyCode = 0x19c
	KEY_DIGITS                   KeyCode = 0x19d
	KEY_TEEN                     KeyCode = 0x19e
	KEY_TWEN                     KeyCode = 0x19f
	KEY_VIDEOPHONE               KeyCode = 0x1a0
	KEY_GAMES                    KeyCode = 0x1a1
	KEY_ZOOMIN                   KeyCode = 0x1a2
	KEY_ZOOMOUT                  KeyCode = 0x1a3
	KEY_ZOOMRESET                KeyCode = 0x1a4
	KEY_WORDPROCESSOR            KeyCode = 0x1a5
	KEY_EDITOR                   KeyCode = 0x1a6
	KEY_SPREADSHEET              KeyCode = 0x1a7
	KEY_GRAPHICSEDITOR           KeyCode = 0x1a8
	KEY_PRESENTATION             KeyCode = 0x1a9
	KEY_DATABASE                 KeyCode = 0x1aa
	KEY_NEWS                     KeyCode = 0x1ab
	KEY_VOICEMAIL                KeyCode = 0x1ac
	KEY_ADDRESSBOOK              KeyCode = 0x1ad
	KEY_MESSENGER                KeyCode = 0x1ae
	KEY_DISPLAYTOGGLE            KeyCode = 0x1af
	KEY_BRIGHTNESS_TOGGLE        KeyCode = 0x1b0
	KEY_SPELLCHECK               KeyCode = 0x1b0
	KEY_LOGOFF                   KeyCode = 0x1b1
	KEY_DOLLAR                   KeyCode = 0x1b2
	KEY_EURO                     KeyCode = 0x1b3
	KEY_FRAMEBACK                KeyCode = 0x1b4
	KEY_FRAMEFORWARD             KeyCode = 0x1b5
	KEY_CONTEXT_MENU             KeyCode = 0x1b6
	KEY_MEDIA_REPEAT             KeyCode = 0x1b7
	KEY_10CHANNELSUP             KeyCode = 0x1b8
	KEY_10CHANNELSDOWN           KeyCode = 0x1b9
	KEY_IMAGES                   KeyCode = 0x1ba
	KEY_DEL_EOL                  KeyCode = 0x1c0
	KEY_DEL_EOS                  KeyCode = 0x1c1
	KEY_INS_LINE                 KeyCode = 0x1c2
	KEY_DEL_LINE                 KeyCode = 0x1c3
	KEY_FN                       KeyCode = 0x1d0
	KEY_FN_ESC                   KeyCode = 0x1d1
	KEY_FN_F1                    KeyCode = 0x1d2
	KEY_FN_F2                    KeyCode = 0x1d3
	KEY_FN_F3                    KeyCode = 0x1d4
	KEY_FN_F4                    KeyCode = 0x1d5
	KEY_FN_F5                    KeyCode = 0x1d6
	KEY_FN_F6                    KeyCode = 0x1d7
	KEY_FN_F7                    KeyCode = 0x1d8
	KEY_FN_F8                    KeyCode = 0x1d9
	KEY_FN_F9                    KeyCode = 0x1da
	KEY_FN_F10                   KeyCode = 0x1db
	KEY_FN_F11                   KeyCode = 0x1dc
	KEY_FN_F12                   KeyCode = 0x1dd
	KEY_FN_1                     KeyCode = 0x1de
	KEY_FN_2                     KeyCode = 0x1df
	KEY_FN_D                     KeyCode = 0x1e0
	KEY_FN_E                     KeyCode = 0x1e1
	KEY_FN_F                     KeyCode = 0x1e2
	KEY_FN_S                     KeyCode = 0x1e3
	KEY_FN_B                     KeyCode = 0x1e4
	KEY_BRL_DOT1                 KeyCode = 0x1f1
	KEY_BRL_DOT2                 KeyCode = 0x1f2
	KEY_BRL_DOT3                 KeyCode = 0x1f3
	KEY_BRL_DOT4                 KeyCode = 0x1f4
	KEY_BRL_DOT5                 KeyCode = 0x1f5
	KEY_BRL_DOT6                 KeyCode = 0x1f6
	KEY_BRL_DOT7                 KeyCode = 0x1f7
	KEY_BRL_DOT8                 KeyCode = 0x1f8
	KEY_BRL_DOT9                 KeyCode = 0x1f9
	KEY_BRL_DOT10                KeyCode = 0x1fa
	KEY_NUMERIC_0                KeyCode = 0x200
	KEY_NUMERIC_1                KeyCode = 0x201
	KEY_NUMERIC_2                KeyCode = 0x202
	KEY_NUMERIC_3                KeyCode = 0x203
	KEY_NUMERIC_4                KeyCode = 0x204
	KEY_NUMERIC_5                KeyCode = 0x205
	KEY_NUMERIC_6                KeyCode = 0x206
	KEY_NUMERIC_7                KeyCode = 0x207
	KEY_NUMERIC_8                KeyCode = 0x208
	KEY_NUMERIC_9                KeyCode = 0x209
	KEY_NUMERIC_STAR             KeyCode = 0x20a
	KEY_NUMERIC_POUND            KeyCode = 0x20b
	KEY_NUMERIC_A                KeyCode = 0x20c
	KEY_NUMERIC_B                KeyCode = 0x20d
	KEY_NUMERIC_C                KeyCode = 0x20e
	KEY_NUMERIC_D                KeyCode = 0x20f
	KEY_CAMERA_FOCUS             KeyCode = 0x210
	KEY_WPS_BUTTON               KeyCode = 0x211
	KEY_TOUCHPAD_TOGGLE          KeyCode = 0x212
	KEY_TOUCHPAD_ON              KeyCode = 0x213
	KEY_TOUCHPAD_OFF             KeyCode = 0x214
	KEY_CAMERA_ZOOMIN            KeyCode = 0x215
	KEY_CAMERA_ZOOMOUT           KeyCode = 0x216
	KEY_CAMERA_UP                KeyCode = 0x217
	KEY_CAMERA_DOWN              KeyCode = 0x218
	KEY_CAMERA_LEFT              KeyCode = 0x219
	KEY_CAMERA_RIGHT             KeyCode = 0x21a
	KEY_ATTENDANT_ON             KeyCode = 0x21b
	KEY_ATTENDANT_OFF            KeyCode = 0x21c
	KEY_ATTENDANT_TOGGLE         KeyCode = 0x21d
	KEY_LIGHTS_TOGGLE            KeyCode = 0x21e
	BTN_DPAD_UP                  KeyCode = 0x220
	BTN_DPAD_DOWN                KeyCode = 0x221
	BTN_DPAD_LEFT                KeyCode = 0x222
	BTN_DPAD_RIGHT               KeyCode = 0x223
	KEY_ALS_TOGGLE               KeyCode = 0x230
	KEY_BUTTONCONFIG             KeyCode = 0x240
	KEY_TASKMANAGER              KeyCode = 0x241
	KEY_JOURNAL                  KeyCode = 0x242
	KEY_CONTROLPANEL             KeyCode = 0x243
	KEY_APPSELECT                KeyCode = 0x244
	KEY_SCREENSAVER              KeyCode = 0x245
	KEY_VOICECOMMAND             KeyCode = 0x246
	KEY_BRIGHTNESS_MIN           KeyCode = 0x250
	KEY_BRIGHTNESS_MAX           KeyCode = 0x251
	KEY_KBDINPUTASSIST_PREV      KeyCode = 0x260
	KEY_KBDINPUTASSIST_NEXT      KeyCode = 0x261
	KEY_KBDINPUTASSIST_PREVGROUP KeyCode = 0x262
	KEY_KBDINPUTASSIST_NEXTGROUP KeyCode = 0x263
	KEY_KBDINPUTASSIST_ACCEPT    KeyCode = 0x264
	KEY_KBDINPUTASSIST_CANCEL    KeyCode = 0x265
	KEY_RIGHT_UP                 KeyCode = 0x266
	KEY_RIGHT_DOWN               KeyCode = 0x267
	KEY_LEFT_UP                  KeyCode = 0x268
	KEY_LEFT_DOWN                KeyCode = 0x269
	KEY_ROOT_MENU                KeyCode = 0x26a
	KEY_MEDIA_TOP_MENU           KeyCode = 0x26b
	KEY_NUMERIC_11               KeyCode = 0x26c
	KEY_NUMERIC_12               KeyCode = 0x26d
	KEY_AUDIO_DESC               KeyCode = 0x26e
	KEY_3D_MODE                  KeyCode = 0x26f
	KEY_NEXT_FAVORITE            KeyCode = 0x270
	KEY_STOP_RECORD              KeyCode = 0x271
	KEY_PAUSE_RECORD             KeyCode = 0x272
	KEY_VOD                      KeyCode = 0x273
	KEY_UNMUTE                   KeyCode = 0x274
	KEY_FASTREVERSE              KeyCode = 0x275
	KEY_SLOWREVERSE              KeyCode = 0x276
	KEY_DATA                     KeyCode = 0x275
	BTN_TRIGGER_HAPPY            KeyCode = 0x2c0
	BTN_TRIGGER_HAPPY1           KeyCode = 0x2c0
	BTN_TRIGGER_HAPPY2           KeyCode = 0x2c1
	BTN_TRIGGER_HAPPY3           KeyCode = 0x2c2
	BTN_TRIGGER_HAPPY4           KeyCode = 0x2c3
	BTN_TRIGGER_HAPPY5           KeyCode = 0x2c4
	BTN_TRIGGER_HAPPY6           KeyCode = 0x2c5
	BTN_TRIGGER_HAPPY7           KeyCode = 0x2c6
	BTN_TRIGGER_HAPPY8           KeyCode = 0x2c7
	BTN_TRIGGER_HAPPY9           KeyCode = 0x2c8
	BTN_TRIGGER_HAPPY10          KeyCode = 0x2c9
	BTN_TRIGGER_HAPPY11          KeyCode = 0x2ca
	BTN_TRIGGER_HAPPY12          KeyCode = 0x2cb
	BTN_TRIGGER_HAPPY13          KeyCode = 0x2cc
	BTN_TRIGGER_HAPPY14          KeyCode = 0x2cd
	BTN_TRIGGER_HAPPY15          KeyCode = 0x2ce
	BTN_TRIGGER_HAPPY16          KeyCode = 0x2cf
	BTN_TRIGGER_HAPPY17          KeyCode = 0x2d0
	BTN_TRIGGER_HAPPY18          KeyCode = 0x2d1
	BTN_TRIGGER_HAPPY19          KeyCode = 0x2d2
	BTN_TRIGGER_HAPPY20          KeyCode = 0x2d3
	BTN_TRIGGER_HAPPY21          KeyCode = 0x2d4
	BTN_TRIGGER_HAPPY22          KeyCode = 0x2d5
	BTN_TRIGGER_HAPPY23          KeyCode = 0x2d6
	BTN_TRIGGER_HAPPY24          KeyCode = 0x2d7
	BTN_TRIGGER_HAPPY25          KeyCode = 0x2d8
	BTN_TRIGGER_HAPPY26          KeyCode = 0x2d9
	BTN_TRIGGER_HAPPY27          KeyCode = 0x2da
	BTN_TRIGGER_HAPPY28          KeyCode = 0x2db
	BTN_TRIGGER_HAPPY29          KeyCode = 0x2dc
	BTN_TRIGGER_HAPPY30          KeyCode = 0x2dd
	BTN_TRIGGER_HAPPY31          KeyCode = 0x2de
	BTN_TRIGGER_HAPPY32          KeyCode = 0x2df
	BTN_TRIGGER_HAPPY33          KeyCode = 0x2e0
	BTN_TRIGGER_HAPPY34          KeyCode = 0x2e1
	BTN_TRIGGER_HAPPY35          KeyCode = 0x2e2
	BTN_TRIGGER_HAPPY36          KeyCode = 0x2e3
	BTN_TRIGGER_HAPPY37          KeyCode = 0x2e4
	BTN_TRIGGER_HAPPY38          KeyCode = 0x2e5
	BTN_TRIGGER_HAPPY39          KeyCode = 0x2e6
	BTN_TRIGGER_HAPPY40          KeyCode = 0x2e7
	KEY_MIN_INTERESTING          KeyCode = 0x2ff
	KEY_MAX                      KeyCode = 0x2ff
)

var ecodes = map[string]KeyCode{
	"KEY_RESERVED":                 KEY_RESERVED,
	"KEY_ESC":                      KEY_ESC,
	"KEY_1":                        KEY_1,
	"KEY_2":                        KEY_2,
	"KEY_3":                        KEY_3,
	"KEY_4":                        KEY_4,
	"KEY_5":                        KEY_5,
	"KEY_6":                        KEY_6,
	"KEY_7":                        KEY_7,
	"KEY_8":                        KEY_8,
	"KEY_9":                        KEY_9,
	"KEY_0":                        KEY_0,
	"KEY_MINUS":                    KEY_MINUS,
	"KEY_EQUAL":                    KEY_EQUAL,
	"KEY_BACKSPACE":                KEY_BACKSPACE,
	"KEY_TAB":                      KEY_TAB,
	"KEY_Q":                        KEY_Q,
	"KEY_W":                        KEY_W,
	"KEY_E":                        KEY_E,
	"KEY_R":                        KEY_R,
	"KEY_T":                        KEY_T,
	"KEY_Y":                        KEY_Y,
	"KEY_U":                        KEY_U,
	"KEY_I":                        KEY_I,
	"KEY_O":                        KEY_O,
	"KEY_P":                        KEY_P,
	"KEY_LEFTBRACE":                KEY_LEFTBRACE,
	"KEY_RIGHTBRACE":               KEY_RIGHTBRACE,
	"KEY_ENTER":                    KEY_ENTER,
	"KEY_LEFTCTRL":                 KEY_LEFTCTRL,
	"KEY_A":                        KEY_A,
	"KEY_S":                        KEY_S,
	"KEY_D":                        KEY_D,
	"KEY_F":                        KEY_F,
	"KEY_G":                        KEY_G,
	"KEY_H":                        KEY_H,
	"KEY_J":                        KEY_J,
	"KEY_K":                        KEY_K,
	"KEY_L":                        KEY_L,
	"KEY_SEMICOLON":                KEY_SEMICOLON,
	"KEY_APOSTROPHE":               KEY_APOSTROPHE,
	"KEY_GRAVE":                    KEY_GRAVE,
	"KEY_LEFTSHIFT":                KEY_LEFTSHIFT,
	"KEY_BACKSLASH":                KEY_BACKSLASH,
	"KEY_Z":                        KEY_Z,
	"KEY_X":                        KEY_X,
	"KEY_C":                        KEY_C,
	"KEY_V":                        KEY_V,
	"KEY_B":                        KEY_B,
	"KEY_N":                        KEY_N,
	"KEY_M":                        KEY_M,
	"KEY_COMMA":                    KEY_COMMA,
	"KEY_DOT":                      KEY_DOT,
	"KEY_SLASH":                    KEY_SLASH,
	"KEY_RIGHTSHIFT":               KEY_RIGHTSHIFT,
	"KEY_KPASTERISK":               KEY_KPASTERISK,
	"KEY_LEFTALT":                  KEY_LEFTALT,
	"KEY_SPACE":                    KEY_SPACE,
	"KEY_CAPSLOCK":                 KEY_CAPSLOCK,
	"KEY_F1":                       KEY_F1,
	"KEY_F2":                       KEY_F2,
	"KEY_F3":                       KEY_F3,
	"KEY_F4":                       KEY_F4,
	"KEY_F5":                       KEY_F5,
	"KEY_F6":                       KEY_F6,
	"KEY_F7":                       KEY_F7,
	"KEY_F8":                       KEY_F8,
	"KEY_F9":                       KEY_F9,
	"KEY_F10":                      KEY_F10,
	"KEY_NUMLOCK":                  KEY_NUMLOCK,
	"KEY_SCROLLLOCK":               KEY_SCROLLLOCK,
	"KEY_KP7":                      KEY_KP7,
	"KEY_KP8":                      KEY_KP8,
	"KEY_KP9":                      KEY_KP9,
	"KEY_KPMINUS":                  KEY_KPMINUS,
	"KEY_KP4":                      KEY_KP4,
	"KEY_KP5":                      KEY_KP5,
	"KEY_KP6":                      KEY_KP6,
	"KEY_KPPLUS":                   KEY_KPPLUS,
	"KEY_KP1":                      KEY_KP1,
	"KEY_KP2":                      KEY_KP2,
	"KEY_KP3":                      KEY_KP3,
	"KEY_KP0":                      KEY_KP0,
	"KEY_KPDOT":                    KEY_KPDOT,
	"KEY_ZENKAKUHANKAKU":           KEY_ZENKAKUHANKAKU,
	"KEY_102ND":                    KEY_102ND,
	"KEY_F11":                      KEY_F11,
	"KEY_F12":                      KEY_F12,
	"KEY_RO":                       KEY_RO,
	"KEY_KATAKANA":                 KEY_KATAKANA,
	"KEY_HIRAGANA":                 KEY_HIRAGANA,
	"KEY_HENKAN":                   KEY_HENKAN,
	"KEY_KATAKANAHIRAGANA":         KEY_KATAKANAHIRAGANA,
	"KEY_MUHENKAN":                 KEY_MUHENKAN,
	"KEY_KPJPCOMMA":                KEY_KPJPCOMMA,
	"KEY_KPENTER":                  KEY_KPENTER,
	"KEY_RIGHTCTRL":                KEY_RIGHTCTRL,
	"KEY_KPSLASH":                  KEY_KPSLASH,
	"KEY_SYSRQ":                    KEY_SYSRQ,
	"KEY_RIGHTALT":                 KEY_RIGHTALT,
	"KEY_LINEFEED":                 KEY_LINEFEED,
	"KEY_HOME":                     KEY_HOME,
	"KEY_UP":                       KEY_UP,
	"KEY_PAGEUP":                   KEY_PAGEUP,
	"KEY_LEFT":                     KEY_LEFT,
	"KEY_RIGHT":                    KEY_RIGHT,
	"KEY_END":                      KEY_END,
	"KEY_DOWN":                     KEY_DOWN,
	"KEY_PAGEDOWN":                 KEY_PAGEDOWN,
	"KEY_INSERT":                   KEY_INSERT,
	"KEY_DELETE":                   KEY_DELETE,
	"KEY_MACRO":                    KEY_MACRO,
	"KEY_MUTE":                     KEY_MUTE,
	"KEY_VOLUMEDOWN":               KEY_VOLUMEDOWN,
	"KEY_VOLUMEUP":                 KEY_VOLUMEUP,
	"KEY_POWER":                    KEY_POWER,
	"KEY_KPEQUAL":                  KEY_KPEQUAL,
	"KEY_KPPLUSMINUS":              KEY_KPPLUSMINUS,
	"KEY_PAUSE":                    KEY_PAUSE,
	"KEY_SCALE":                    KEY_SCALE,
	"KEY_KPCOMMA":                  KEY_KPCOMMA,
	"KEY_HANGEUL":                  KEY_HANGEUL,
	"KEY_HANGUEL":                  KEY_HANGUEL,
	"KEY_HANJA":                    KEY_HANJA,
	"KEY_YEN":                      KEY_YEN,
	"KEY_LEFTMETA":                 KEY_LEFTMETA,
	"KEY_RIGHTMETA":                KEY_RIGHTMETA,
	"KEY_COMPOSE":                  KEY_COMPOSE,
	"KEY_STOP":                     KEY_STOP,
	"KEY_AGAIN":                    KEY_AGAIN,
	"KEY_PROPS":                    KEY_PROPS,
	"KEY_UNDO":                     KEY_UNDO,
	"KEY_FRONT":                    KEY_FRONT,
	"KEY_COPY":                     KEY_COPY,
	"KEY_OPEN":                     KEY_OPEN,
	"KEY_PASTE":                    KEY_PASTE,
	"KEY_FIND":                     KEY_FIND,
	"KEY_CUT":                      KEY_CUT,
	"KEY_HELP":                     KEY_HELP,
	"KEY_MENU":                     KEY_MENU,
	"KEY_CALC":                     KEY_CALC,
	"KEY_SETUP":                    KEY_SETUP,
	"KEY_SLEEP":                    KEY_SLEEP,
	"KEY_WAKEUP":                   KEY_WAKEUP,
	"KEY_FILE":                     KEY_FILE,
	"KEY_SENDFILE":                 KEY_SENDFILE,
	"KEY_DELETEFILE":               KEY_DELETEFILE,
	"KEY_XFER":                     KEY_XFER,
	"KEY_PROG1":                    KEY_PROG1,
	"KEY_PROG2":                    KEY_PROG2,
	"KEY_WWW":                      KEY_WWW,
	"KEY_MSDOS":                    KEY_MSDOS,
	"KEY_COFFEE":                   KEY_COFFEE,
	"KEY_SCREENLOCK":               KEY_SCREENLOCK,
	"KEY_ROTATE_DISPLAY":           KEY_ROTATE_DISPLAY,
	"KEY_DIRECTION":                KEY_DIRECTION,
	"KEY_CYCLEWINDOWS":             KEY_CYCLEWINDOWS,
	"KEY_MAIL":                     KEY_MAIL,
	"KEY_BOOKMARKS":                KEY_BOOKMARKS,
	"KEY_COMPUTER":                 KEY_COMPUTER,
	"KEY_BACK":                     KEY_BACK,
	"KEY_FORWARD":                  KEY_FORWARD,
	"KEY_CLOSECD":                  KEY_CLOSECD,
	"KEY_EJECTCD":                  KEY_EJECTCD,
	"KEY_EJECTCLOSECD":             KEY_EJECTCLOSECD,
	"KEY_NEXTSONG":                 KEY_NEXTSONG,
	"KEY_PLAYPAUSE":                KEY_PLAYPAUSE,
	"KEY_PREVIOUSSONG":             KEY_PREVIOUSSONG,
	"KEY_STOPCD":                   KEY_STOPCD,
	"KEY_RECORD":                   KEY_RECORD,
	"KEY_REWIND":                   KEY_REWIND,
	"KEY_PHONE":                    KEY_PHONE,
	"KEY_ISO":                      KEY_ISO,
	"KEY_CONFIG":                   KEY_CONFIG,
	"KEY_HOMEPAGE":                 KEY_HOMEPAGE,
	"KEY_REFRESH":                  KEY_REFRESH,
	"KEY_EXIT":                     KEY_EXIT,
	"KEY_MOVE":                     KEY_MOVE,
	"KEY_EDIT":                     KEY_EDIT,
	"KEY_SCROLLUP":                 KEY_SCROLLUP,
	"KEY_SCROLLDOWN":               KEY_SCROLLDOWN,
	"KEY_KPLEFTPAREN":              KEY_KPLEFTPAREN,
	"KEY_KPRIGHTPAREN":             KEY_KPRIGHTPAREN,
	"KEY_NEW":                      KEY_NEW,
	"KEY_REDO":                     KEY_REDO,
	"KEY_F13":                      KEY_F13,
	"KEY_F14":                      KEY_F14,
	"KEY_F15":                      KEY_F15,
	"KEY_F16":                      KEY_F16,
	"KEY_F17":                      KEY_F17,
	"KEY_F18":                      KEY_F18,
	"KEY_F19":                      KEY_F19,
	"KEY_F20":                      KEY_F20,
	"KEY_F21":                      KEY_F21,
	"KEY_F22":                      KEY_F22,
	"KEY_F23":                      KEY_F23,
	"KEY_F24":                      KEY_F24,
	"KEY_PLAYCD":                   KEY_PLAYCD,
	"KEY_PAUSECD":                  KEY_PAUSECD,
	"KEY_PROG3":                    KEY_PROG3,
	"KEY_PROG4":                    KEY_PROG4,
	"KEY_DASHBOARD":                KEY_DASHBOARD,
	"KEY_SUSPEND":                  KEY_SUSPEND,
	"KEY_CLOSE":                    KEY_CLOSE,
	"KEY_PLAY":                     KEY_PLAY,
	"KEY_FASTFORWARD":              KEY_FASTFORWARD,
	"KEY_BASSBOOST":                KEY_BASSBOOST,
	"KEY_PRINT":                    KEY_PRINT,
	"KEY_HP":                       KEY_HP,
	"KEY_CAMERA":                   KEY_CAMERA,
	"KEY_SOUND":                    KEY_SOUND,
	"KEY_QUESTION":                 KEY_QUESTION,
	"KEY_EMAIL":                    KEY_EMAIL,
	"KEY_CHAT":                     KEY_CHAT,
	"KEY_SEARCH":                   KEY_SEARCH,
	"KEY_CONNECT":                  KEY_CONNECT,
	"KEY_FINANCE":                  KEY_FINANCE,
	"KEY_SPORT":                    KEY_SPORT,
	"KEY_SHOP":                     KEY_SHOP,
	"KEY_ALTERASE":                 KEY_ALTERASE,
	"KEY_CANCEL":                   KEY_CANCEL,
	"KEY_BRIGHTNESSDOWN":           KEY_BRIGHTNESSDOWN,
	"KEY_BRIGHTNESSUP":             KEY_BRIGHTNESSUP,
	"KEY_MEDIA":                    KEY_MEDIA,
	"KEY_SWITCHVIDEOMODE":          KEY_SWITCHVIDEOMODE,
	"KEY_KBDILLUMTOGGLE":           KEY_KBDILLUMTOGGLE,
	"KEY_KBDILLUMDOWN":             KEY_KBDILLUMDOWN,
	"KEY_KBDILLUMUP":               KEY_KBDILLUMUP,
	"KEY_SEND":                     KEY_SEND,
	"KEY_REPLY":                    KEY_REPLY,
	"KEY_FORWARDMAIL":              KEY_FORWARDMAIL,
	"KEY_SAVE":                     KEY_SAVE,
	"KEY_DOCUMENTS":                KEY_DOCUMENTS,
	"KEY_BATTERY":                  KEY_BATTERY,
	"KEY_BLUETOOTH":                KEY_BLUETOOTH,
	"KEY_WLAN":                     KEY_WLAN,
	"KEY_UWB":                      KEY_UWB,
	"KEY_UNKNOWN":                  KEY_UNKNOWN,
	"KEY_VIDEO_NEXT":               KEY_VIDEO_NEXT,
	"KEY_VIDEO_PREV":               KEY_VIDEO_PREV,
	"KEY_BRIGHTNESS_CYCLE":         KEY_BRIGHTNESS_CYCLE,
	"KEY_BRIGHTNESS_AUTO":          KEY_BRIGHTNESS_AUTO,
	"KEY_DISPLAY_OFF":              KEY_DISPLAY_OFF,
	"KEY_WWAN":                     KEY_WWAN,
	"KEY_RFKILL":                   KEY_RFKILL,
	"KEY_MICMUTE":                  KEY_MICMUTE,
	"BTN_MISC":                     BTN_MISC,
	"BTN_0":                        BTN_0,
	"BTN_1":                        BTN_1,
	"BTN_2":                        BTN_2,
	"BTN_3":                        BTN_3,
	"BTN_4":                        BTN_4,
	"BTN_5":                        BTN_5,
	"BTN_6":                        BTN_6,
	"BTN_7":                        BTN_7,
	"BTN_8":                        BTN_8,
	"BTN_9":                        BTN_9,
	"BTN_MOUSE":                    BTN_MOUSE,
	"BTN_LEFT":                     BTN_LEFT,
	"BTN_RIGHT":                    BTN_RIGHT,
	"BTN_MIDDLE":                   BTN_MIDDLE,
	"BTN_SIDE":                     BTN_SIDE,
	"BTN_EXTRA":                    BTN_EXTRA,
	"BTN_FORWARD":                  BTN_FORWARD,
	"BTN_BACK":                     BTN_BACK,
	"BTN_TASK":                     BTN_TASK,
	"BTN_JOYSTICK":                 BTN_JOYSTICK,
	"BTN_TRIGGER":                  BTN_TRIGGER,
	"BTN_THUMB":                    BTN_THUMB,
	"BTN_THUMB2":                   BTN_THUMB2,
	"BTN_TOP":                      BTN_TOP,
	"BTN_TOP2":                     BTN_TOP2,
	"BTN_PINKIE":                   BTN_PINKIE,
	"BTN_BASE":                     BTN_BASE,
	"BTN_BASE2":                    BTN_BASE2,
	"BTN_BASE3":                    BTN_BASE3,
	"BTN_BASE4":                    BTN_BASE4,
	"BTN_BASE5":                    BTN_BASE5,
	"BTN_BASE6":                    BTN_BASE6,
	"BTN_DEAD":                     BTN_DEAD,
	"BTN_GAMEPAD":                  BTN_GAMEPAD,
	"BTN_SOUTH":                    BTN_SOUTH,
	"BTN_A":                        BTN_A,
	"BTN_EAST":                     BTN_EAST,
	"BTN_B":                        BTN_B,
	"BTN_C":                        BTN_C,
	"BTN_NORTH":                    BTN_NORTH,
	"BTN_X":                        BTN_X,
	"BTN_WEST":                     BTN_WEST,
	"BTN_Y":                        BTN_Y,
	"BTN_Z":                        BTN_Z,
	"BTN_TL":                       BTN_TL,
	"BTN_TR":                       BTN_TR,
	"BTN_TL2":                      BTN_TL2,
	"BTN_TR2":                      BTN_TR2,
	"BTN_SELECT":                   BTN_SELECT,
	"BTN_START":                    BTN_START,
	"BTN_MODE":                     BTN_MODE,
	"BTN_THUMBL":                   BTN_THUMBL,
	"BTN_THUMBR":                   BTN_THUMBR,
	"BTN_DIGI":                     BTN_DIGI,
	"BTN_TOOL_PEN":                 BTN_TOOL_PEN,
	"BTN_TOOL_RUBBER":              BTN_TOOL_RUBBER,
	"BTN_TOOL_BRUSH":               BTN_TOOL_BRUSH,
	"BTN_TOOL_PENCIL":              BTN_TOOL_PENCIL,
	"BTN_TOOL_AIRBRUSH":            BTN_TOOL_AIRBRUSH,
	"BTN_TOOL_FINGER":              BTN_TOOL_FINGER,
	"BTN_TOOL_MOUSE":               BTN_TOOL_MOUSE,
	"BTN_TOOL_LENS":                BTN_TOOL_LENS,
	"BTN_TOOL_QUINTTAP":            BTN_TOOL_QUINTTAP,
	"BTN_TOUCH":                    BTN_TOUCH,
	"BTN_STYLUS":                   BTN_STYLUS,
	"BTN_STYLUS2":                  BTN_STYLUS2,
	"BTN_TOOL_DOUBLETAP":           BTN_TOOL_DOUBLETAP,
	"BTN_TOOL_TRIPLETAP":           BTN_TOOL_TRIPLETAP,
	"BTN_TOOL_QUADTAP":             BTN_TOOL_QUADTAP,
	"BTN_WHEEL":                    BTN_WHEEL,
	"BTN_GEAR_DOWN":                BTN_GEAR_DOWN,
	"BTN_GEAR_UP":                  BTN_GEAR_UP,
	"KEY_OK":                       KEY_OK,
	"KEY_SELECT":                   KEY_SELECT,
	"KEY_GOTO":                     KEY_GOTO,
	"KEY_CLEAR":                    KEY_CLEAR,
	"KEY_POWER2":                   KEY_POWER2,
	"KEY_OPTION":                   KEY_OPTION,
	"KEY_INFO":                     KEY_INFO,
	"KEY_TIME":                     KEY_TIME,
	"KEY_VENDOR":                   KEY_VENDOR,
	"KEY_ARCHIVE":                  KEY_ARCHIVE,
	"KEY_PROGRAM":                  KEY_PROGRAM,
	"KEY_CHANNEL":                  KEY_CHANNEL,
	"KEY_FAVORITES":                KEY_FAVORITES,
	"KEY_EPG":                      KEY_EPG,
	"KEY_PVR":                      KEY_PVR,
	"KEY_MHP":                      KEY_MHP,
	"KEY_LANGUAGE":                 KEY_LANGUAGE,
	"KEY_TITLE":                    KEY_TITLE,
	"KEY_SUBTITLE":                 KEY_SUBTITLE,
	"KEY_ANGLE":                    KEY_ANGLE,
	"KEY_ZOOM":                     KEY_ZOOM,
	"KEY_MODE":                     KEY_MODE,
	"KEY_KEYBOARD":                 KEY_KEYBOARD,
	"KEY_SCREEN":                   KEY_SCREEN,
	"KEY_PC":                       KEY_PC,
	"KEY_TV":                       KEY_TV,
	"KEY_TV2":                      KEY_TV2,
	"KEY_VCR":                      KEY_VCR,
	"KEY_VCR2":                     KEY_VCR2,
	"KEY_SAT":                      KEY_SAT,
	"KEY_SAT2":                     KEY_SAT2,
	"KEY_CD":                       KEY_CD,
	"KEY_TAPE":                     KEY_TAPE,
	"KEY_RADIO":                    KEY_RADIO,
	"KEY_TUNER":                    KEY_TUNER,
	"KEY_PLAYER":                   KEY_PLAYER,
	"KEY_TEXT":                     KEY_TEXT,
	"KEY_DVD":                      KEY_DVD,
	"KEY_AUX":                      KEY_AUX,
	"KEY_MP3":                      KEY_MP3,
	"KEY_AUDIO":                    KEY_AUDIO,
	"KEY_VIDEO":                    KEY_VIDEO,
	"KEY_DIRECTORY":                KEY_DIRECTORY,
	"KEY_LIST":                     KEY_LIST,
	"KEY_MEMO":                     KEY_MEMO,
	"KEY_CALENDAR":                 KEY_CALENDAR,
	"KEY_RED":                      KEY_RED,
	"KEY_GREEN":                    KEY_GREEN,
	"KEY_YELLOW":                   KEY_YELLOW,
	"KEY_BLUE":                     KEY_BLUE,
	"KEY_CHANNELUP":                KEY_CHANNELUP,
	"KEY_CHANNELDOWN":              KEY_CHANNELDOWN,
	"KEY_FIRST":                    KEY_FIRST,
	"KEY_LAST":                     KEY_LAST,
	"KEY_AB":                       KEY_AB,
	"KEY_NEXT":                     KEY_NEXT,
	"KEY_RESTART":                  KEY_RESTART,
	"KEY_SLOW":                     KEY_SLOW,
	"KEY_SHUFFLE":                  KEY_SHUFFLE,
	"KEY_BREAK":                    KEY_BREAK,
	"KEY_PREVIOUS":                 KEY_PREVIOUS,
	"KEY_DIGITS":                   KEY_DIGITS,
	"KEY_TEEN":                     KEY_TEEN,
	"KEY_TWEN":                     KEY_TWEN,
	"KEY_VIDEOPHONE":               KEY_VIDEOPHONE,
	"KEY_GAMES":                    KEY_GAMES,
	"KEY_ZOOMIN":                   KEY_ZOOMIN,
	"KEY_ZOOMOUT":                  KEY_ZOOMOUT,
	"KEY_ZOOMRESET":                KEY_ZOOMRESET,
	"KEY_WORDPROCESSOR":            KEY_WORDPROCESSOR,
	"KEY_EDITOR":                   KEY_EDITOR,
	"KEY_SPREADSHEET":              KEY_SPREADSHEET,
	"KEY_GRAPHICSEDITOR":           KEY_GRAPHICSEDITOR,
	"KEY_PRESENTATION":             KEY_PRESENTATION,
	"KEY_DATABASE":                 KEY_DATABASE,
	"KEY_NEWS":                     KEY_NEWS,
	"KEY_VOICEMAIL":                KEY_VOICEMAIL,
	"KEY_ADDRESSBOOK":              KEY_ADDRESSBOOK,
	"KEY_MESSENGER":                KEY_MESSENGER,
	"KEY_DISPLAYTOGGLE":            KEY_DISPLAYTOGGLE,
	"KEY_BRIGHTNESS_TOGGLE":        KEY_BRIGHTNESS_TOGGLE,
	"KEY_SPELLCHECK":               KEY_SPELLCHECK,
	"KEY_LOGOFF":                   KEY_LOGOFF,
	"KEY_DOLLAR":                   KEY_DOLLAR,
	"KEY_EURO":                     KEY_EURO,
	"KEY_FRAMEBACK":                KEY_FRAMEBACK,
	"KEY_FRAMEFORWARD":             KEY_FRAMEFORWARD,
	"KEY_CONTEXT_MENU":             KEY_CONTEXT_MENU,
	"KEY_MEDIA_REPEAT":             KEY_MEDIA_REPEAT,
	"KEY_10CHANNELSUP":             KEY_10CHANNELSUP,
	"KEY_10CHANNELSDOWN":           KEY_10CHANNELSDOWN,
	"KEY_IMAGES":                   KEY_IMAGES,
	"KEY_DEL_EOL":                  KEY_DEL_EOL,
	"KEY_DEL_EOS":                  KEY_DEL_EOS,
	"KEY_INS_LINE":                 KEY_INS_LINE,
	"KEY_DEL_LINE":                 KEY_DEL_LINE,
	"KEY_FN":                       KEY_FN,
	"KEY_FN_ESC":                   KEY_FN_ESC,
	"KEY_FN_F1":                    KEY_FN_F1,
	"KEY_FN_F2":                    KEY_FN_F2,
	"KEY_FN_F3":                    KEY_FN_F3,
	"KEY_FN_F4":                    KEY_FN_F4,
	"KEY_FN_F5":                    KEY_FN_F5,
	"KEY_FN_F6":                    KEY_FN_F6,
	"KEY_FN_F7":                    KEY_FN_F7,
	"KEY_FN_F8":                    KEY_FN_F8,
	"KEY_FN_F9":                    KEY_FN_F9,
	"KEY_FN_F10":                   KEY_FN_F10,
	"KEY_FN_F11":                   KEY_FN_F11,
	"KEY_FN_F12":                   KEY_FN_F12,
	"KEY_FN_1":                     KEY_FN_1,
	"KEY_FN_2":                     KEY_FN_2,
	"KEY_FN_D":                     KEY_FN_D,
	"KEY_FN_E":                     KEY_FN_E,
	"KEY_FN_F":                     KEY_FN_F,
	"KEY_FN_S":                     KEY_FN_S,
	"KEY_FN_B":                     KEY_FN_B,
	"KEY_BRL_DOT1":                 KEY_BRL_DOT1,
	"KEY_BRL_DOT2":                 KEY_BRL_DOT2,
	"KEY_BRL_DOT3":                 KEY_BRL_DOT3,
	"KEY_BRL_DOT4":                 KEY_BRL_DOT4,
	"KEY_BRL_DOT5":                 KEY_BRL_DOT5,
	"KEY_BRL_DOT6":                 KEY_BRL_DOT6,
	"KEY_BRL_DOT7":                 KEY_BRL_DOT7,
	"KEY_BRL_DOT8":                 KEY_BRL_DOT8,
	"KEY_BRL_DOT9":                 KEY_BRL_DOT9,
	"KEY_BRL_DOT10":                KEY_BRL_DOT10,
	"KEY_NUMERIC_0":                KEY_NUMERIC_0,
	"KEY_NUMERIC_1":                KEY_NUMERIC_1,
	"KEY_NUMERIC_2":                KEY_NUMERIC_2,
	"KEY_NUMERIC_3":                KEY_NUMERIC_3,
	"KEY_NUMERIC_4":                KEY_NUMERIC_4,
	"KEY_NUMERIC_5":                KEY_NUMERIC_5,
	"KEY_NUMERIC_6":                KEY_NUMERIC_6,
	"KEY_NUMERIC_7":                KEY_NUMERIC_7,
	"KEY_NUMERIC_8":                KEY_NUMERIC_8,
	"KEY_NUMERIC_9":                KEY_NUMERIC_9,
	"KEY_NUMERIC_STAR":             KEY_NUMERIC_STAR,
	"KEY_NUMERIC_POUND":            KEY_NUMERIC_POUND,
	"KEY_NUMERIC_A":                KEY_NUMERIC_A,
	"KEY_NUMERIC_B":                KEY_NUMERIC_B,
	"KEY_NUMERIC_C":                KEY_NUMERIC_C,
	"KEY_NUMERIC_D":                KEY_NUMERIC_D,
	"KEY_CAMERA_FOCUS":             KEY_CAMERA_FOCUS,
	"KEY_WPS_BUTTON":               KEY_WPS_BUTTON,
	"KEY_TOUCHPAD_TOGGLE":          KEY_TOUCHPAD_TOGGLE,
	"KEY_TOUCHPAD_ON":              KEY_TOUCHPAD_ON,
	"KEY_TOUCHPAD_OFF":             KEY_TOUCHPAD_OFF,
	"KEY_CAMERA_ZOOMIN":            KEY_CAMERA_ZOOMIN,
	"KEY_CAMERA_ZOOMOUT":           KEY_CAMERA_ZOOMOUT,
	"KEY_CAMERA_UP":                KEY_CAMERA_UP,
	"KEY_CAMERA_DOWN":              KEY_CAMERA_DOWN,
	"KEY_CAMERA_LEFT":              KEY_CAMERA_LEFT,
	"KEY_CAMERA_RIGHT":             KEY_CAMERA_RIGHT,
	"KEY_ATTENDANT_ON":             KEY_ATTENDANT_ON,
	"KEY_ATTENDANT_OFF":            KEY_ATTENDANT_OFF,
	"KEY_ATTENDANT_TOGGLE":         KEY_ATTENDANT_TOGGLE,
	"KEY_LIGHTS_TOGGLE":            KEY_LIGHTS_TOGGLE,
	"BTN_DPAD_UP":                  BTN_DPAD_UP,
	"BTN_DPAD_DOWN":                BTN_DPAD_DOWN,
	"BTN_DPAD_LEFT":                BTN_DPAD_LEFT,
	"BTN_DPAD_RIGHT":               BTN_DPAD_RIGHT,
	"KEY_ALS_TOGGLE":               KEY_ALS_TOGGLE,
	"KEY_BUTTONCONFIG":             KEY_BUTTONCONFIG,
	"KEY_TASKMANAGER":              KEY_TASKMANAGER,
	"KEY_JOURNAL":                  KEY_JOURNAL,
	"KEY_CONTROLPANEL":             KEY_CONTROLPANEL,
	"KEY_APPSELECT":                KEY_APPSELECT,
	"KEY_SCREENSAVER":              KEY_SCREENSAVER,
	"KEY_VOICECOMMAND":             KEY_VOICECOMMAND,
	"KEY_BRIGHTNESS_MIN":           KEY_BRIGHTNESS_MIN,
	"KEY_BRIGHTNESS_MAX":           KEY_BRIGHTNESS_MAX,
	"KEY_KBDINPUTASSIST_PREV":      KEY_KBDINPUTASSIST_PREV,
	"KEY_KBDINPUTASSIST_NEXT":      KEY_KBDINPUTASSIST_NEXT,
	"KEY_KBDINPUTASSIST_PREVGROUP": KEY_KBDINPUTASSIST_PREVGROUP,
	"KEY_KBDINPUTASSIST_NEXTGROUP": KEY_KBDINPUTASSIST_NEXTGROUP,
	"KEY_KBDINPUTASSIST_ACCEPT":    KEY_KBDINPUTASSIST_ACCEPT,
	"KEY_KBDINPUTASSIST_CANCEL":    KEY_KBDINPUTASSIST_CANCEL,
	"KEY_RIGHT_UP":                 KEY_RIGHT_UP,
	"KEY_RIGHT_DOWN":               KEY_RIGHT_DOWN,
	"KEY_LEFT_UP":                  KEY_LEFT_UP,
	"KEY_LEFT_DOWN":                KEY_LEFT_DOWN,
	"KEY_ROOT_MENU":                KEY_ROOT_MENU,
	"KEY_MEDIA_TOP_MENU":           KEY_MEDIA_TOP_MENU,
	"KEY_NUMERIC_11":               KEY_NUMERIC_11,
	"KEY_NUMERIC_12":               KEY_NUMERIC_12,
	"KEY_AUDIO_DESC":               KEY_AUDIO_DESC,
	"KEY_3D_MODE":                  KEY_3D_MODE,
	"KEY_NEXT_FAVORITE":            KEY_NEXT_FAVORITE,
	"KEY_STOP_RECORD":              KEY_STOP_RECORD,
	"KEY_PAUSE_RECORD":             KEY_PAUSE_RECORD,
	"KEY_VOD":                      KEY_VOD,
	"KEY_UNMUTE":                   KEY_UNMUTE,
	"KEY_FASTREVERSE":              KEY_FASTREVERSE,
	"KEY_SLOWREVERSE":              KEY_SLOWREVERSE,
	"KEY_DATA":                     KEY_DATA,
	"BTN_TRIGGER_HAPPY":            BTN_TRIGGER_HAPPY,
	"BTN_TRIGGER_HAPPY1":           BTN_TRIGGER_HAPPY1,
	"BTN_TRIGGER_HAPPY2":           BTN_TRIGGER_HAPPY2,
	"BTN_TRIGGER_HAPPY3":           BTN_TRIGGER_HAPPY3,
	"BTN_TRIGGER_HAPPY4":           BTN_TRIGGER_HAPPY4,
	"BTN_TRIGGER_HAPPY5":           BTN_TRIGGER_HAPPY5,
	"BTN_TRIGGER_HAPPY6":           BTN_TRIGGER_HAPPY6,
	"BTN_TRIGGER_HAPPY7":           BTN_TRIGGER_HAPPY7,
	"BTN_TRIGGER_HAPPY8":           BTN_TRIGGER_HAPPY8,
	"BTN_TRIGGER_HAPPY9":           BTN_TRIGGER_HAPPY9,
	"BTN_TRIGGER_HAPPY10":          BTN_TRIGGER_HAPPY10,
	"BTN_TRIGGER_HAPPY11":          BTN_TRIGGER_HAPPY11,
	"BTN_TRIGGER_HAPPY12":          BTN_TRIGGER_HAPPY12,
	"BTN_TRIGGER_HAPPY13":          BTN_TRIGGER_HAPPY13,
	"BTN_TRIGGER_HAPPY14":          BTN_TRIGGER_HAPPY14,
	"BTN_TRIGGER_HAPPY15":          BTN_TRIGGER_HAPPY15,
	"BTN_TRIGGER_HAPPY16":          BTN_TRIGGER_HAPPY16,
	"BTN_TRIGGER_HAPPY17":          BTN_TRIGGER_HAPPY17,
	"BTN_TRIGGER_HAPPY18":          BTN_TRIGGER_HAPPY18,
	"BTN_TRIGGER_HAPPY19":          BTN_TRIGGER_HAPPY19,
	"BTN_TRIGGER_HAPPY20":          BTN_TRIGGER_HAPPY20,
	"BTN_TRIGGER_HAPPY21":          BTN_TRIGGER_HAPPY21,
	"BTN_TRIGGER_HAPPY22":          BTN_TRIGGER_HAPPY22,
	"BTN_TRIGGER_HAPPY23":          BTN_TRIGGER_HAPPY23,
	"BTN_TRIGGER_HAPPY24":          BTN_TRIGGER_HAPPY24,
	"BTN_TRIGGER_HAPPY25":          BTN_TRIGGER_HAPPY25,
	"BTN_TRIGGER_HAPPY26":          BTN_TRIGGER_HAPPY26,
	"BTN_TRIGGER_HAPPY27":          BTN_TRIGGER_HAPPY27,
	"BTN_TRIGGER_HAPPY28":          BTN_TRIGGER_HAPPY28,
	"BTN_TRIGGER_HAPPY29":          BTN_TRIGGER_HAPPY29,
	"BTN_TRIGGER_HAPPY30":          BTN_TRIGGER_HAPPY30,
	"BTN_TRIGGER_HAPPY31":          BTN_TRIGGER_HAPPY31,
	"BTN_TRIGGER_HAPPY32":          BTN_TRIGGER_HAPPY32,
	"BTN_TRIGGER_HAPPY33":          BTN_TRIGGER_HAPPY33,
	"BTN_TRIGGER_HAPPY34":          BTN_TRIGGER_HAPPY34,
	"BTN_TRIGGER_HAPPY35":          BTN_TRIGGER_HAPPY35,
	"BTN_TRIGGER_HAPPY36":          BTN_TRIGGER_HAPPY36,
	"BTN_TRIGGER_HAPPY37":          BTN_TRIGGER_HAPPY37,
	"BTN_TRIGGER_HAPPY38":          BTN_TRIGGER_HAPPY38,
	"BTN_TRIGGER_HAPPY39":          BTN_TRIGGER_HAPPY39,
	"BTN_TRIGGER_HAPPY40":          BTN_TRIGGER_HAPPY40,
	"KEY_MIN_INTERESTING":          KEY_MIN_INTERESTING,
	"KEY_MAX":                      KEY_MAX,
}
