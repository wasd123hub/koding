package models

import (
	"errors"
	// "fmt"
	"github.com/jinzhu/gorm"
	"github.com/koding/bongo"
	"time"
)

type Notification struct {
	// unique identifier of Notification
	Id int64 `json:"id"`

	// notification recipient account id
	AccountId int64 `json:"accountId" sql:"NOT NULL"`

	// notification content foreign key
	NotificationContentId int64 `json:"notificationContentId" sql:"NOT NULL"`

	// glanced information
	Glanced bool `json:"glanced" sql:"NOT NULL"`

	// last notifier addition time
	UpdatedAt time.Time `json:"updatedAt" sql:"NOT NULL"`
}

const (
	Notification_TYPE_SUBSCRIBE   = "subscribe"
	Notification_TYPE_UNSUBSCRIBE = "unsubscribe"
)

func (n *Notification) GetId() int64 {
	return n.Id
}

func (n Notification) TableName() string {
	return "api.notification"
}

func NewNotification() *Notification {
	return &Notification{}
}

func (n *Notification) One(q *bongo.Query) error {
	return bongo.B.One(n, n, q)
}

func (n *Notification) Create() error {
	s := map[string]interface{}{
		"account_id":              n.AccountId,
		"notification_content_id": n.NotificationContentId,
	}
	q := bongo.NewQS(s)
	if err := n.One(q); err != nil {
		if err != gorm.RecordNotFound {
			return err
		}

		return bongo.B.Create(n)
	}

	n.Glanced = false

	return bongo.B.Update(n)
}

func (n *Notification) List(q *Query) (*NotificationResponse, error) {
	if q.Limit == 0 {
		return nil, errors.New("limit cannot be zero")
	}
	response := &NotificationResponse{}
	result, err := n.getDecoratedList(q)
	if err != nil {
		return response, err
	}

	response.Notifications = result
	response.UnreadCount = getUnreadNotificationCount(result)

	return response, nil
}

func (n *Notification) Some(data interface{}, q *bongo.Query) error {

	return bongo.B.Some(n, data, q)
}

func (n *Notification) fetchByAccountId(q *Query) ([]Notification, error) {
	var notifications []Notification
	query := &bongo.Query{
		Selector: map[string]interface{}{
			"account_id": q.AccountId,
		},
		Sort: map[string]string{
			"updated_at": "desc",
		},
		Pagination: bongo.Pagination{
			Limit: q.Limit,
		},
	}
	if err := bongo.B.Some(n, &notifications, query); err != nil {
		return nil, err
	}

	return notifications, nil
}

// getDecoratedList fetches notifications of the given user and decorates it with
// notification activity actors
func (n *Notification) getDecoratedList(q *Query) ([]NotificationContainer, error) {
	result := make([]NotificationContainer, 0)

	nList, err := n.fetchByAccountId(q)
	if err != nil {
		return nil, err
	}
	// fetch all notification content relationships
	contentIds := deductContentIds(nList)

	nc := NewNotificationContent()
	ncMap, err := nc.FetchMapByIds(contentIds)
	if err != nil {
		return nil, err
	}

	na := NewNotificationActivity()
	naMap, err := na.FetchMapByContentIds(contentIds)
	if err != nil {
		return nil, err
	}

	for _, n := range nList {
		nc := ncMap[n.NotificationContentId]
		na := naMap[n.NotificationContentId]
		container := n.buildNotificationContainer(q.AccountId, &nc, na)
		result = append(result, container)
	}

	return result, nil
}

func (n *Notification) buildNotificationContainer(actorId int64, nc *NotificationContent, na []NotificationActivity) NotificationContainer {
	ct, err := CreateNotificationContentType(nc.TypeConstant)
	if err != nil {
		return NotificationContainer{}
	}

	ct.SetTargetId(nc.TargetId)
	ct.SetListerId(actorId)
	ac, err := ct.FetchActors(na)
	if err != nil {
		return NotificationContainer{}
	}
	return NotificationContainer{
		TargetId:              nc.TargetId,
		TypeConstant:          nc.TypeConstant,
		UpdatedAt:             n.UpdatedAt,
		Glanced:               n.Glanced,
		NotificationContentId: nc.Id,
		LatestActors:          ac.LatestActors,
		ActorCount:            ac.Count,
	}
}

func deductContentIds(nList []Notification) []int64 {
	notificationContentIds := make([]int64, 0)
	for _, n := range nList {
		notificationContentIds = append(notificationContentIds, n.NotificationContentId)
	}

	return notificationContentIds
}

func (n *Notification) FetchContent() (*NotificationContent, error) {
	nc := NewNotificationContent()
	if err := nc.ById(n.NotificationContentId); err != nil {
		return nil, err
	}

	return nc, nil
}

// func (n *Notification) Follow(a *NotificationActivity) error {
// 	// a.TypeConstant = NotificationContent_TYPE_FOLLOW
// 	// create NotificationActivity
// 	if err := a.Create(); err != nil {
// 		return err
// 	}

// 	fn := NewFollowNotification()
// 	fn.NotifierId = a.ActorId
// 	// fn.TargetId = a.TargetId

// 	return CreateNotificationContent(fn)
// }

// func (n *Notification) JoinGroup(a *NotificationActivity, admins []int64) error {
// 	// a.TypeConstant = NotificationContent_TYPE_JOIN

// 	return n.interactGroup(a, admins)
// }

// func (n *Notification) LeaveGroup(a *NotificationActivity, admins []int64) error {
// 	// a.TypeConstant = NotificationContent_TYPE_LEAVE

// 	return n.interactGroup(a, admins)
// }

// func (n *Notification) interactGroup(a *NotificationActivity, admins []int64) error {
// 	gn := NewGroupNotification(a.TypeConstant)
// 	gn.NotifierId = a.ActorId
// 	gn.TargetId = a.TargetId
// 	gn.Admins = admins

// 	if err := a.Create(); err != nil {
// 		return err
// 	}

// 	return CreateNotificationContent(gn)
// }

func (n *Notification) AfterCreate() {
	bongo.B.AfterCreate(n)
}

func (n *Notification) AfterUpdate() {
	bongo.B.AfterUpdate(n)
}

func (n *Notification) AfterDelete() {
	bongo.B.AfterDelete(n)
}

func (n *Notification) Glance() error {
	selector := map[string]interface{}{
		"glanced":    false,
		"account_id": n.AccountId,
	}

	set := map[string]interface{}{
		"glanced": true,
	}

	return bongo.B.UpdateMulti(n, selector, set)
}

func getUnreadNotificationCount(notificationList []NotificationContainer) int {
	unreadCount := 0
	for _, nc := range notificationList {
		if !nc.Glanced {
			unreadCount++
		}
	}

	return unreadCount
}
