package handlers_test

/*func TestDeletePermission(t *testing.T) {
	// Create a new controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Define the test cases
	tests := []struct {
		name           string
		permissionID   string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:         "Without UUID param",
			permissionID: "",
			mockSetup: func() {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false).Times(1),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), "admin.permissions.delete").Return(false).Times(0),
				)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "Without permission",
			permissionID: "valid-uuid",
			mockSetup: func() {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), "admin.permissions.delete").Return(false),
				)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:         "With delete error",
			permissionID: "valid-uuid",
			mockSetup: func() {
				m := mocks.NewMockUtils(ctrl)
				gomock.InOrder(
					m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true),
					m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), "admin.permissions.delete").Return(true),
				)
				handlers.ReplaceGlobals(m)
				permissions.ReplaceGlobals(mocks.NewPermissionsRepository(mocks.PermissionsRepository{Error: errors.New("error")}))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:         "With success",
			permissionID: "valid-uuid",
			mockSetup: func() {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), "admin.permissions.delete").Return(true)
				handlers.ReplaceGlobals(m)
				permissions.ReplaceGlobals(mocks.NewPermissionsRepository(mocks.PermissionsRepository{}))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new recorder and request
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/permissions/"+tt.permissionID, nil)

			// Apply mocks
			tt.mockSetup()

			// Call the function
			handlers.DeletePermission(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.Status)
		})
	}
}
*/
