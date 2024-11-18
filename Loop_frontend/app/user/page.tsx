"use client";

import { useEffect, useState } from "react";
import { getUserInfo } from "./actions";
import { Card, CardBody, CardHeader } from "@nextui-org/card";
import { Avatar,AvatarIcon } from "@nextui-org/avatar";
import { Divider } from "@nextui-org/divider";
import { Chip } from "@nextui-org/chip";
import { subheading, heading } from "@/components/ui/primitives";
import { useAuthStore } from "@/lib/auth/authStore";
import { UserType } from "../types";

export default function UserPage() {
  const [user, setUser] = useState<UserType | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const queryParams = new URLSearchParams(window.location.search);
  const user_id = queryParams.get("id");

  const refresh_token = useAuthStore((state) => state.refresh_token);

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        if (refresh_token && user_id){
          const userInfo = await getUserInfo(refresh_token,user_id);
          setUser(userInfo);
        }
      } catch (err) {
        console.error("Error fetching user data:", err);
        setError("Failed to fetch user details.");
      } finally {
        setLoading(false);
      }
    };

    fetchUserData();
  }, [user_id, refresh_token]);

  if (loading) {
    return <div>Loading user details...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!user) {
    return <div>No user data found.</div>;
  }

  return (  
    <div className="w-11/12 mx-auto px-4 py-8 p-12">
      <Card className="mb-8 p-10">
        <CardBody>
          <div className="flex flex-col md:flex-row gap-8 items-start">
            <div className="flex flex-col items-center space-y-4">
              <Avatar
              icon={<AvatarIcon />}
              classNames={{
                base: "bg-gradient-to-br from-[#00B4DB] to-[#0083B0]",
                icon: "text-black/80",
              }}
              className="w-48 h-48"
              alt={user.name}
              />
              <div className={heading({ size: "lg" })}>{user.name}</div>
              <Chip color="primary" variant="flat">User</Chip>
            </div>

            <div className="flex-1 space-y-6">
              <div>
                <p className={subheading({ size: "sm" })}>{user.bio}</p>
              </div>

              <Divider />

              <div className="grid grid-cols-2 gap-4">
                <Card>
                  <CardBody>
                    <div className={subheading({ size: "sm" })}>Location</div>
                    <p className="text-foreground-600">{user.location}</p>
                  </CardBody>
                </Card>

                <Card>
                  <CardBody>
                    <div className={subheading({ size: "sm" })}>Joined</div>
                    <p className="text-foreground-600">
                      {new Date(user.created_at).toLocaleDateString("en-US", {
                        year: "numeric",
                        month: "long",
                        day: "numeric",
                      })}
                    </p>
                  </CardBody>
                </Card>
              </div>
            </div>
          </div>
        </CardBody>
      </Card>

      <Card className="p-6">
        <CardHeader>
          <div className={heading({ size: "md" })}>Projects</div>
        </CardHeader>
        <CardBody>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {user.projects?.map((project) => (
              <a href={`/projectpage?id=${project.project_id}`}>
              <Card key={project.project_id} isPressable>
                <CardBody>
                  <img
                    src="https://www.liquidplanner.com/wp-content/uploads/2019/04/HiRes-17.jpg"
                    alt={project.title}
                    className="w-full h-48 object-cover rounded-lg mb-4"
                  />
                  <div className={subheading({ size: "sm" })}>{project.title}</div>
                  <p className="text-foreground-600 text-sm">{project.description}</p>
                </CardBody>
              </Card>
              </a>
            ))}
          </div>
        </CardBody>
      </Card>
    </div>
  );
}
