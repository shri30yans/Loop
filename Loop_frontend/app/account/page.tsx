"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { Spinner } from "@nextui-org/spinner";
import { useAuthStore } from "../../lib/auth/authStore";
import { Button } from "@nextui-org/button";
import { Card, CardFooter, CardBody, CardHeader } from "@nextui-org/card";
import {
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  useDisclosure,
} from "@nextui-org/modal";

export default function AccountPage() {
  const [isLoading, setIsLoading] = useState(false);
  const { isOpen, onOpen, onClose } = useDisclosure();

  const user_id = useAuthStore.getState().user_id;
  const logout = useAuthStore((state) => state.logout);
  const router = useRouter();

  const handleLogout = async () => {
    try {
      setIsLoading(true);
      await logout();
      router.push("/");
    } catch (error) {
      console.error("Logout failed:", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-8">Account Settings</h1>

        <div className="grid gap-6 md:grid-cols-2">
          <Card className="p-6">
            <h2 className="text-xl font-semibold mb-4">Profile</h2>
            <div className="space-y-4">
              <a
                href={`/profile?id=${user_id}`}
                className="text-blue-600 hover:underline block"
              >
                View Public Profile
              </a>
              {/* <Button
                color="primary"
                variant="flat"
                onClick={() => router.push("/account/edit")}
              >
                Edit Profile
              </Button> */}
            </div>
          </Card>

          <Card className="p-6">
            <h2 className="text-xl font-semibold mb-4">Security</h2>
            <div className="flex gap-2">
              <Button
                color="primary"
                variant="flat"
                onClick={() => router.push("/auth/edit_password")}
              >
                Change Password
              </Button>
              <Button
                color="danger"
                variant="flat"
                onClick={onOpen}
                disabled={isLoading}
              >
                {isLoading ? <Spinner size="sm" /> : "Logout"}
              </Button>
            </div>
          </Card>
        </div>

        <Modal isOpen={isOpen} onClose={onClose} size="md">
          <ModalContent>
            {(onClose) => (
              <>
                <ModalHeader>Confirm Logout</ModalHeader>
                <ModalBody>Are you sure you want to logout?</ModalBody>
                <ModalFooter>
                  <Button color="danger" variant="light" onPress={onClose}>
                    Cancel
                  </Button>
                  <Button
                    color="primary"
                    onPress={handleLogout}
                    isLoading={isLoading}
                  >
                    Logout
                  </Button>
                </ModalFooter>
              </>
            )}
          </ModalContent>
        </Modal>
      </div>
  );
}
