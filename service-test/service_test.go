package service_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/t-ksn/user-service/src/service"
)

var _ = Describe("User service:", func() {
	var username string
	var password string
	BeforeEach(func() {
		username = safeGenString()
		password = safeGenString()
	})

	Describe("User create account", func() {
		Context("With password more then 4 chars", func() {
			It("Should get success confirmation", func() {
				err := userServiceClient.Register(username, password)

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("With password less then 4 chars", func() {
			It("Should get 'Minimum password length 4' error", func() {
				err := userServiceClient.Register(username, "12")

				Expect(err).Should(Equal(service.ErrPasswordLessThen4Chars))
			})
		})
		Context("With duplicated user name", func() {
			It("Should get 'User name already exist' error", func() {
				err := userServiceClient.Register(username, password)
				Expect(err).ShouldNot(HaveOccurred())

				err = userServiceClient.Register(username, "new password")
				Expect(err).Should(Equal(service.ErrDuplicateName))
			})
		})
		Context("With empty user name", func() {
			It("Should get 'User name is empty' error", func() {
				err := userServiceClient.Register("", password)

				Expect(err).Should(Equal(service.ErrUserNameIsEmpty))
			})

		})
	})
	Describe("User has account", func() {
		var password string
		BeforeEach(func() {
			password = safeGenString()

			err := userServiceClient.Register(username, password)
			Expect(err).ShouldNot(HaveOccurred(), "Seccess registration")
		})

		Describe("and try to login", func() {

			Context("With wrong password", func() {
				It("Should get 'user name or passowrd incorrect' error", func() {
					_, err := userServiceClient.SignIn(username, "wrong_password")
					Expect(err).Should(Equal(service.ErrUserNameOrPasswordIncorrect))
				})
			})
			Context("With wrong username", func() {
				It("should get 'user name or passowrd incorrect' error", func() {
					_, err := userServiceClient.SignIn("wrong_username", password)
					Expect(err).Should(Equal(service.ErrUserNameOrPasswordIncorrect))
				})
			})
			It("should get access token", func() {
				_, err := userServiceClient.SignIn(username, password)
				Expect(err).ShouldNot(HaveOccurred(), "Seccess login")
			})
		})
	})
})
