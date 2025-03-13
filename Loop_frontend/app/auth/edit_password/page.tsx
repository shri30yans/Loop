'use client'
import { FormEvent, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from "@nextui-org/button";
import { Card, CardFooter, CardBody, CardHeader } from "@nextui-org/card";
import { Input } from '@nextui-org/input';
import { Divider } from "@nextui-org/divider";
import { updatePassword } from './actions';
import { useAuthStore } from '../../../lib/auth/authStore';

export default function EditPasswordPage() {
  const router = useRouter();
  const [error, setError] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);
  const access_token = useAuthStore((state) => state.access_token);

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setError('');
    setIsLoading(true);

    const formData = new FormData(event.currentTarget as HTMLFormElement);
    const currentPassword = formData.get('currentPassword') as string;
    const newPassword = formData.get('newPassword') as string;
    const confirmPassword = formData.get('confirmPassword') as string;

    if (newPassword !== confirmPassword) {
      setError('Passwords do not match');
      setIsLoading(false);
      return;
    }

    try {
      if (access_token) {
        await updatePassword(access_token, currentPassword, newPassword);
        router.push('/account');
      }
    } catch (err) {
      setError('Failed to update password');
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div className="fixed inset-0 flex items-center justify-center">
      <Card className="w-1/4">
        <CardHeader className="flex items-center justify-center pt-6">
          <h1 className="text-2xl font-bold">Change Password</h1>
        </CardHeader>
        <CardBody className="p-6">
          {error && (
            <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
              {error}
            </div>
          )}
          <form onSubmit={handleSubmit} className="space-y-4">
            <Input
              label="Current Password"
              name="currentPassword"
              type="password"
              variant="bordered"
              isRequired
              className="w-full"
              isDisabled={isLoading}
            />
            <Input
              label="New Password"
              name="newPassword"
              type="password"
              variant="bordered"
              isRequired
              className="w-full"
              isDisabled={isLoading}
            />
            <Input
              label="Confirm New Password"
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
              {isLoading ? 'Updating...' : 'Update Password'}
            </Button>
          </form>
        </CardBody>
        <CardFooter className="flex flex-col gap-2 px-6 pb-6">
          <Divider className="my-4" />
          <p className="text-center text-sm text-gray-600">
            Want to go back?{' '}
            <a href="/profile" className="text-blue-600 hover:underline">
              Return to Profile
            </a>
          </p>
        </CardFooter>
      </Card>
    </div>
  );
}
