// LoginPage.tsx
"use client";
import { FormEvent, useState } from "react";
import { login } from "./actions";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { NetworkError, AuthError, ERROR_MESSAGES } from "@/utils/errors";
import { 
  Button,
  Card, 
  CardFooter, 
  CardBody, 
  CardHeader,
  Input,
  Divider 
} from "@nextui-org/react";
import { useAuth } from "../../../lib/auth/authStore";
import Link from "next/link";

export default function LoginPage() {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setError("");
    setIsLoading(true);

    const formData = new FormData(event.currentTarget as HTMLFormElement);
    const email = formData.get("email") as string;
    const password = formData.get("password") as string;

    try {
      const response = await login(email, password);
      if (response) {
        await useAuth.getState().login(email, password);
        const authState = useAuth.getState();

        if (authState.isAuthenticated) {
          toast.success("Login successful!");
          router.push("/");
        } else {
          const authError = authState.error;
          const errorMessage = authError instanceof AuthError ? 
            ERROR_MESSAGES[authError.code] : 
            authError instanceof NetworkError ? 
              authError.message : 
              "Invalid email or password";
          
          toast.error(errorMessage);
          setError(errorMessage);
        }
      }
    } catch (error: any) {
      console.error("Login error:", error);
      console.error("Error details:", {
        name: error.name,
        message: error.message,
        stack: error.stack,
      });

      const errorMessage = error instanceof AuthError ? 
        ERROR_MESSAGES[error.code] :
        error instanceof NetworkError ? 
          error.message : 
          "Login failed. Please try again.";
      
      toast.error(errorMessage);
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div className="fixed inset-0 flex items-center justify-center">
      <Card className="w-1/4">
        <CardHeader className="flex items-center justify-center pt-6">
          <h1 className="text-2xl font-bold">Login</h1>
        </CardHeader>
        <CardBody className="p-6">
          {error && (
            <div className={`mb-4 p-3 rounded flex flex-col gap-2 ${
              error.includes("Server is currently unavailable") 
                ? "bg-yellow-100 border border-yellow-400 text-yellow-700"
                : error.includes("check your internet connection")
                ? "bg-orange-100 border border-orange-400 text-orange-700"
                : "bg-red-100 border border-red-400 text-red-700"
            }`}>
              <div className="font-medium">{error}</div>
              {error.includes("No account found") && (
                <div className="text-sm">
                  New to Loop? <Link href="/auth/register" className="text-blue-600 hover:underline">Create an account</Link>
                </div>
              )}
              {error.includes("Incorrect password") && (
                <div className="text-sm">
                  Forgot your password? <Link href="/auth/reset-password" className="text-blue-600 hover:underline">Reset it here</Link>
                </div>
              )}
              {error.includes("Server is currently unavailable") && (
                <div className="text-sm">
                  Our team has been notified. Please try again in a few minutes.
                </div>
              )}
            </div>
          )}
          <form onSubmit={handleSubmit} className="space-y-4">
            <Input
              label="Email"
              name="email"
              type="email"
              variant="bordered"
              isRequired
              className="w-full"
              isDisabled={isLoading}
            />
            <Input
              label="Password"
              name="password"
              type="password"
              variant="bordered"
              isRequired
              className="w-full"
              isDisabled={isLoading}
            />
            <Button
              type="submit"
              color="primary"
              className="w-full"
              isLoading={isLoading}
            >
              {isLoading ? "Logging in..." : "Login"}
            </Button>
          </form>
        </CardBody>
        <CardFooter className="flex flex-col gap-2 px-6 pb-6">
          <Divider className="my-4" />
          <p className="text-center text-sm text-gray-600">
            Don&apos;t have an account?{" "}
            <Link
              href="/auth/register"
              className="text-blue-600 hover:underline"
            >
              Register here
            </Link>
          </p>
        </CardFooter>
      </Card>
    </div>
  );
}
