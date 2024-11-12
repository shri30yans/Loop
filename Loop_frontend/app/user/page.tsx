'use client'
import { useAuthStore } from "@/lib/auth/authStore";
import { Button } from "@nextui-org/button";
import { useRouter } from "next/navigation";

export default function Home() {
  const logout = useAuthStore((state) => state.logout);
  const router = useRouter();

  const handleLogout = () => {
    logout();
    router.push("/"); // Redirect to home page after logout
  };

  return (
    <section className="flex flex-col items-center justify-center gap-4">
      <h1>USER PAGE</h1>
      <Button 
        color="danger"
        variant="flat"
        onClick={handleLogout}
      >
        Logout
      </Button>
    </section>
  );
}