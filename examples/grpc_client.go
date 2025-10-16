package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/nurefendi/auth-service-using-golang/grpc/pb"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create auth service client
	client := pb.NewAuthServiceClient(conn)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Example 1: Register a new user
	log.Println("=== Testing User Registration ===")
	registerReq := &pb.RegisterRequest{
		FullName: "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
		Gender:   1,
	}

	registerResp, err := client.Register(ctx, registerReq)
	if err != nil {
		log.Printf("Registration failed: %v", err)
	} else {
		log.Printf("Registration response: %+v", registerResp)
	}

	// Example 2: Login
	log.Println("\n=== Testing User Login ===")
	loginReq := &pb.LoginRequest{
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	loginResp, err := client.Login(ctx, loginReq)
	if err != nil {
		log.Printf("Login failed: %v", err)
	} else {
		log.Printf("Login response: %+v", loginResp)
		if loginResp.Success {
			log.Printf("Access Token: %s", loginResp.AccessToken)
			log.Printf("User Info: %+v", loginResp.User)

			// Example 3: Get user profile with JWT token
			log.Println("\n=== Testing Get User Profile ===")
			// Add Authorization header
			md := metadata.Pairs("authorization", "Bearer "+loginResp.AccessToken)
			ctxWithAuth := metadata.NewOutgoingContext(ctx, md)

			profileReq := &pb.GetUserProfileRequest{}
			profileResp, err := client.GetUserProfile(ctxWithAuth, profileReq)
			if err != nil {
				log.Printf("Get profile failed: %v", err)
			} else {
				log.Printf("Profile response: %+v", profileResp)
			}

			// Example 4: Check access
			log.Println("\n=== Testing Check Access ===")
			checkAccessReq := &pb.CheckAccessRequest{
				Path:   "/api/users",
				Method: "GET",
			}

			checkAccessResp, err := client.CheckAccess(ctxWithAuth, checkAccessReq)
			if err != nil {
				log.Printf("Check access failed: %v", err)
			} else {
				log.Printf("Check access response: %+v", checkAccessResp)
			}

			// Example 5: Get user functions
			log.Println("\n=== Testing Get User Functions ===")
			functionsReq := &pb.GetUserFunctionsRequest{}
			functionsResp, err := client.GetUserFunctions(ctxWithAuth, functionsReq)
			if err != nil {
				log.Printf("Get functions failed: %v", err)
			} else {
				log.Printf("Functions response: %+v", functionsResp)
			}

			// Example 6: Refresh token
			log.Println("\n=== Testing Refresh Token ===")
			refreshReq := &pb.RefreshTokenRequest{
				RefreshToken: loginResp.RefreshToken,
			}

			refreshResp, err := client.RefreshToken(ctx, refreshReq)
			if err != nil {
				log.Printf("Refresh token failed: %v", err)
			} else {
				log.Printf("Refresh token response: %+v", refreshResp)
			}

			// Example 7: Logout
			log.Println("\n=== Testing Logout ===")
			logoutReq := &pb.LogoutRequest{
				RefreshToken: loginResp.RefreshToken,
			}

			logoutResp, err := client.Logout(ctxWithAuth, logoutReq)
			if err != nil {
				log.Printf("Logout failed: %v", err)
			} else {
				log.Printf("Logout response: %+v", logoutResp)
			}
		}
	}

	// Example 8: Change password
	log.Println("\n=== Testing Change Password ===")
	changePasswordReq := &pb.ChangePasswordRequest{
		CurrentPassword: "password123",
		NewPassword:     "newpassword123",
	}

	changePasswordResp, err := client.ChangePassword(ctx, changePasswordReq)
	if err != nil {
		log.Printf("Change password failed: %v", err)
	} else {
		log.Printf("Change password response: %+v", changePasswordResp)
	}

	log.Println("\n=== gRPC Client Testing Completed ===")
}
