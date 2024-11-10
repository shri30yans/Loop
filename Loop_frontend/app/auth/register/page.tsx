// RegisterPage.tsx
'use client'
import { FormEvent, useState } from 'react';
import { register } from './actions';
import { useRouter } from 'next/navigation';
import { Button } from "@nextui-org/button";
import { Card, CardFooter, CardBody, CardHeader } from "@nextui-org/card";
import { Input } from '@nextui-org/input';
import { Divider } from "@nextui-org/divider";
import Link from 'next/link';

export default function RegisterPage() {
  const router = useRouter();
  const [error, setError] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setError('');
    setIsLoading(true);

    const formData = new FormData(event.currentTarget as HTMLFormElement);
    const username = formData.get('username') as string;
    const email = formData.get('email') as string;
    const password = formData.get('password') as string;
    const confirmPassword = formData.get('confirmPassword') as string;

    if (password !== confirmPassword) {
      setError('Passwords do not match');
      setIsLoading(false);
      return;
    }

    try {
      const data = await register(username,email, password);
      router.push('/auth/login');
    } catch (error) {
      setError('Registration failed. Please try again.');
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div className="fixed inset-0 flex items-center justify-center">
      <Card className="w-full max-w-md">
        <CardHeader className="flex items-center justify-center pt-6">
          <h1 className="text-2xl font-bold">Register</h1>
        </CardHeader>
        <CardBody className="p-6">
          {error && (
            <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
              {error}
            </div>
          )}
          <form onSubmit={handleSubmit} className="space-y-4">
          <Input
              label="Username"
              name="username"
              variant="bordered"
              isRequired
              className="w-full"
              isDisabled={isLoading}
            />
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
            <Input
              label="Confirm Password"
              name="confirmPassword"
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
              {isLoading ? 'Registering...' : 'Register'}
            </Button>
          </form>
        </CardBody>
        <CardFooter className="flex flex-col gap-2 px-6 pb-6">
          <Divider className="my-4" />
          <p className="text-center text-sm text-gray-600">
            Already have an account?{' '}
            <Link href="/auth/login" className="text-blue-600 hover:underline">
              Login here
            </Link>
          </p>
        </CardFooter>
      </Card>
    </div>
  );
}