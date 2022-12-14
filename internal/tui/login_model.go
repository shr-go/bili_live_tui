package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shr-go/bili_live_tui/api"
	"github.com/shr-go/bili_live_tui/internal/live_room"
	"github.com/shr-go/bili_live_tui/pkg/logging"
	"net/http"
	"os"
	"time"
)

type loginStep uint8

const (
	loginStepConfirmLogin loginStep = iota
	loginStepWaitLogin
	loginStepLoginNeedRefresh
	loginStepLoginSuccess
	loginStepDone
)

type loginModel struct {
	step        loginStep
	client      *http.Client
	loginData   *api.QRCodeLoginData
	room        *api.LiveRoom
	cookies     string
	chooseLogin bool
	localCookie bool
	quit        bool
}

func newLoginModel(client *http.Client) loginModel {
	return loginModel{
		step:        loginStepConfirmLogin,
		client:      client,
		loginData:   nil,
		room:        nil,
		cookies:     "",
		chooseLogin: true,
		localCookie: false,
		quit:        false,
	}
}

type waitScanMsg struct{}

type TickMsg time.Time

func (m *loginModel) loadLoginData() tea.Msg {
	loginData, err := live_room.QRCodeLogin(m.client)
	if err != nil {
		logging.Fatalf("loadLoginData failed, err=%v", err)
	}
	m.loginData = loginData
	return waitScanMsg{}
}

func tickEvery() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m *loginModel) pollLoginStatus() tea.Msg {
	cookies, err := live_room.PollLogin(m.client, m.loginData)
	if err != nil {
		logging.Fatalf("pollLoginStatus failed, err=%v", err)
	}
	switch m.loginData.Status {
	case api.QRLoginExpired:
		m.step = loginStepLoginNeedRefresh
	case api.QRLoginSuccess:
		m.step = loginStepLoginSuccess
		m.cookies = cookies
	}
	return m.step
}

func (m *loginModel) enterRoom() tea.Msg {
	if m.chooseLogin && !m.localCookie {
		if !live_room.CheckCookieValid(m.client, m.cookies) {
			logging.Fatalf("PrepareEnterRoom cookies check failed, program exit")
		}
		os.WriteFile("COOKIE.DAT", []byte(m.cookies), 0o660)
	}

	if room, err := live_room.AuthAndConnect(m.client, LiveConfig.RoomID); err != nil {
		logging.Fatalf("AuthAndConnect failed, err=%v", err)
	} else {
		m.room = room
	}
	m.step = loginStepDone
	return m.step
}

func (m *loginModel) Init() tea.Cmd {
	if m.step == loginStepLoginSuccess {
		return m.enterRoom
	}
	return nil
}

func (m *loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quit = true
			return m, tea.Quit
		}
		switch m.step {
		case loginStepConfirmLogin:
			switch msg.String() {
			case "tab":
				m.chooseLogin = !m.chooseLogin
			case "left":
				m.chooseLogin = true
			case "right":
				m.chooseLogin = false
			case "enter", " ":
				if m.chooseLogin {
					m.step = loginStepWaitLogin
					return m, m.loadLoginData
				} else {
					m.step = loginStepLoginSuccess
					return m, m.enterRoom
				}
			}
		case loginStepLoginNeedRefresh:
			if msg.String() == "enter" || msg.String() == " " {
				m.step = loginStepWaitLogin
				m.loginData = nil
				return m, m.loadLoginData
			}
		}
		return m, nil
	case waitScanMsg:
		return m, tickEvery()
	case TickMsg:
		return m, m.pollLoginStatus
	case loginStep:
		if msg == loginStepWaitLogin {
			return m, tickEvery()
		} else if msg == loginStepLoginSuccess {
			return m, m.enterRoom
		} else if msg == loginStepDone {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *loginModel) View() string {
	switch m.step {
	case loginStepConfirmLogin:
		var loginButton, cancelButton string
		if m.chooseLogin {
			loginButton = activeButtonStyle.Render("????????????")
			cancelButton = buttonStyle.Render("??????")
		} else {
			loginButton = buttonStyle.Render("????????????")
			cancelButton = activeButtonStyle.Render("??????")
		}

		question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).
			Render("???????????????????????????????????????")
		buttons := lipgloss.JoinHorizontal(lipgloss.Top, loginButton, "  ", cancelButton)
		ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)
		dialog := lipgloss.Place(windowWidth, windowHeight,
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(ui),
			lipgloss.WithWhitespaceForeground(subtle),
		)
		return dialog
	case loginStepWaitLogin:
		if m.loginData != nil {
			tips := "?????????????????????????????????????????????login.png????????????"
			if m.loginData.Status == api.QRLoginNotConfirm {
				tips = "???????????????????????????????????????"
			}
			ui := lipgloss.JoinVertical(lipgloss.Center, tips, m.loginData.QRString)
			dialogBoxStyleCopy := dialogBoxStyle.Copy().Padding(0, 0)
			return lipgloss.Place(windowWidth, windowHeight,
				lipgloss.Center, lipgloss.Center,
				dialogBoxStyleCopy.Render(ui),
				lipgloss.WithWhitespaceForeground(subtle),
			)
		}
	case loginStepLoginNeedRefresh:
		question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).
			Render("???????????????????????????????????????")
		confirmButton := activeButtonStyle.Render("??????")
		ui := lipgloss.JoinVertical(lipgloss.Center, question, confirmButton)
		return lipgloss.Place(windowWidth, windowHeight,
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(ui),
			lipgloss.WithWhitespaceForeground(subtle),
		)
	case loginStepLoginSuccess:
		str := "????????????????????????????????????"
		if !m.chooseLogin {
			str = "?????????????????????????????????????????????"
		}
		text := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(str)
		dialog := lipgloss.Place(windowWidth, windowHeight,
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(text),
			lipgloss.WithWhitespaceForeground(subtle),
		)
		return dialog
	}
	return ""
}
