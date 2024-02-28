package model

// func BugFromGrpcMessage(r *coreService.Bug) *Bug {
// 	return &Bug{
// 		ID:                r.GetBugID(),
// 		Name:              newPointerStr(r.GetName()),
// 		Description:       newPointerStr(r.GetDescription()),
// 		CreatedForRelease: &Release{ID: r.GetCreateForReleaseID()},
// 		SolvedInRelease:   &Release{ID: r.GetSolvedInReleaseID()},
// 		CreatedAt:         newPointerStr(r.GetCreatedAt().AsTime().Format("02.01.2006")),
// 		UpdatedAt:         newPointerStr(r.GetUpdatedAt().AsTime().Format("02.01.2006")),
// 		State:             getStateFromGrpc(r.GetState()),
// 		Files:             FilesFromGrpcMessage(r.GetFiles()),
// 	}
// }

// func BugListResponseFromGrpcMessage(r *coreService.GetBugListRes) *BugsResponse {
// 	var bugs []*Bug
// 	for _, v := range r.GetBugs() {
// 		bugs = append(bugs, BugFromGrpcMessage(v))
// 	}

// 	return &BugsResponse{
// 		TotalCount: newPointerInt(int(r.GetTotalCount())),
// 		TotalPages: newPointerInt(int(r.GetTotalPages())),
// 		Page:       newPointerInt(int(r.GetPage())),
// 		Size:       newPointerInt(int(r.GetSize())),
// 		HasMore:    newPointerBool(r.GetHasMore()),
// 		Bugs:       bugs,
// 	}
// }

// func SearchResponseFromGrpcMessage(r *coreService.SearchBugRes) *BugsResponse {
// 	var bugs []*Bug
// 	for _, v := range r.GetBugs() {
// 		bugs = append(bugs, BugFromGrpcMessage(v))
// 	}

// 	return &BugsResponse{
// 		TotalCount: newPointerInt(int(r.GetTotalCount())),
// 		TotalPages: newPointerInt(int(r.GetTotalPages())),
// 		Page:       newPointerInt(int(r.GetPage())),
// 		Size:       newPointerInt(int(r.GetSize())),
// 		HasMore:    newPointerBool(r.GetHasMore()),
// 		Bugs:       bugs,
// 	}
// }

// func GetFavouritesResponseFromGrpcMessage(r *coreService.GetFavouritesRes) *BugsResponse {
// 	var bugs []*Bug
// 	for _, v := range r.GetBugs() {
// 		bugs = append(bugs, BugFromGrpcMessage(v))
// 	}

// 	return &BugsResponse{
// 		TotalCount: newPointerInt(int(r.GetTotalCount())),
// 		TotalPages: newPointerInt(int(r.GetTotalPages())),
// 		Page:       newPointerInt(int(r.GetPage())),
// 		Size:       newPointerInt(int(r.GetSize())),
// 		HasMore:    newPointerBool(r.GetHasMore()),
// 		Bugs:       bugs,
// 	}
// }

// func CommentFromGrpcMessage(r *coreService.Comment) *Comment {
// 	return &Comment{
// 		Comment: newPointerStr(r.GetComment()),
// 		User:    &User{ID: newPointerStr(r.GetUserUUID())},
// 		Date:    newPointerStr(r.GetDate().AsTime().Format("02.01.2006")),
// 		Files:   FilesFromGrpcMessage(r.GetFiles()),
// 	}
// }

// func CommentsResponseFromGrpcMessage(r *coreService.GetCommentsRes) *CommentsResponse {
// 	var comments []*Comment
// 	for _, v := range r.GetComments() {
// 		comments = append(comments, CommentFromGrpcMessage(v))
// 	}

// 	return &CommentsResponse{
// 		Result:     CommentFromGrpcMessage(r.GetResult()),
// 		Workaround: CommentFromGrpcMessage(r.GetWorkaround()),
// 		Comments:   comments,
// 	}
// }

// func FilesFromGrpcMessage(files []*coreService.File) []*File {
// 	var res []*File
// 	for _, f := range files {
// 		res = append(res, FileFromGrpcMessage(f))
// 	}

// 	return res
// }

// func FileFromGrpcMessage(f *coreService.File) *File {
// 	return &File{
// 		FileID:      newPointerStr(f.FileID),
// 		PreviewData: newPointerStr(string(f.Data)),
// 		FullFileID:  newPointerStr(f.FullFileID),
// 	}
// }

// func VotesResponseFromGrpcMessage(r *coreService.GetVotesRes) *VotesResponse {
// 	return &VotesResponse{
// 		Like:    newPointerInt64(r.GetLikes()),
// 		Dislike: newPointerInt64(r.GetDislikes()),
// 	}
// }

// func newPointerStr(in string) *string {
// 	res := in
// 	return &res
// }

// func newPointerInt(in int) *int {
// 	res := in
// 	return &res
// }

// func newPointerInt64(in int64) *int {
// 	var res int
// 	res = int(in)
// 	return &res
// }

// func newPointerBool(in bool) *bool {
// 	res := in
// 	return &res
// }

// func getStateFromGrpc(s coreService.BugStates) *BugState {
// 	res := BugStateNew
// 	if s == coreService.BugStates_NEW {
// 		res = BugStateNew
// 	} else if s == coreService.BugStates_INWORK {
// 		res = BugStateInwork
// 	} else if s == coreService.BugStates_DONE {
// 		res = BugStateDone
// 	} else if s == coreService.BugStates_REJECTED {
// 		res = BugStateRejected
// 	}
// 	return &res
// }

// func GetBugStateToGrpc(s *BugState) coreService.BugStates {
// 	res := coreService.BugStates_ALL
// 	if s == nil {
// 		res = coreService.BugStates_ALL
// 	} else if *s == BugStateNew {
// 		res = coreService.BugStates_NEW
// 	} else if *s == BugStateInwork {
// 		res = coreService.BugStates_INWORK
// 	} else if *s == BugStateDone {
// 		res = coreService.BugStates_DONE
// 	} else if *s == BugStateRejected {
// 		res = coreService.BugStates_REJECTED
// 	}
// 	return res
// }

// func GetCommentStateToGrpc(s CommentState) coreService.CommentStates {
// 	res := coreService.CommentStates_COMMENT
// 	if s == CommentStateComment {
// 		res = coreService.CommentStates_COMMENT
// 	} else if s == CommentStateWorkaround {
// 		res = coreService.CommentStates_WORKAROUND
// 	} else if s == CommentStateResult {
// 		res = coreService.CommentStates_RESULT
// 	}
// 	return res
// }

// func GetVoteStateToGrpc(s VoteState) coreService.VoteStates {
// 	res := coreService.VoteStates_LIKE
// 	if s == VoteStateLike {
// 		res = coreService.VoteStates_LIKE
// 	} else if s == VoteStateDislike {
// 		res = coreService.VoteStates_DISLIKE
// 	}
// 	return res
// }
