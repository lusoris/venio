'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { apiClient, User } from '@/lib/api';

export default function DashboardPage() {
  const { user, logout } = useAuth();
  const router = useRouter();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!user) {
      router.push('/login');
      return;
    }

    const fetchUsers = async () => {
      try {
        setLoading(true);
        const data = await apiClient.getUsers({ page: 1, limit: 10 });
        setUsers(data.users);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load users');
      } finally {
        setLoading(false);
      }
    };

    fetchUsers();
  }, [user, router]);

  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 to-gray-800">
      {/* Header */}
      <header className="bg-gray-800 border-b border-gray-700">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <h1 className="text-2xl font-bold text-white">Venio</h1>
            </div>
            <div className="flex items-center space-x-4">
              <div className="text-sm text-gray-400">
                <span className="font-medium text-white">{user.username}</span>
                <span className="mx-2">•</span>
                <span>{user.email}</span>
              </div>
              <button
                onClick={logout}
                className="px-4 py-2 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-lg transition-colors"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h2 className="text-3xl font-bold text-white mb-2">Dashboard</h2>
          <p className="text-gray-400">Welcome back, {user.username}!</p>
        </div>

        {/* User Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
            <h3 className="text-gray-400 text-sm font-medium mb-2">Your Roles</h3>
            <div className="flex flex-wrap gap-2">
              {user.roles?.map((role) => (
                <span
                  key={role}
                  className="px-3 py-1 bg-blue-600/20 text-blue-400 rounded-full text-sm font-medium"
                >
                  {role}
                </span>
              ))}
              {(!user.roles || user.roles.length === 0) && (
                <span className="text-gray-500 text-sm">No roles assigned</span>
              )}
            </div>
          </div>

          <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
            <h3 className="text-gray-400 text-sm font-medium mb-2">Account Status</h3>
            <p className="text-2xl font-bold text-green-400">Active</p>
          </div>

          <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
            <h3 className="text-gray-400 text-sm font-medium mb-2">Total Users</h3>
            <p className="text-2xl font-bold text-white">{users.length}</p>
          </div>
        </div>

        {/* Users List */}
        <div className="bg-gray-800 rounded-lg border border-gray-700">
          <div className="px-6 py-4 border-b border-gray-700">
            <h3 className="text-xl font-semibold text-white">Users</h3>
          </div>
          <div className="p-6">
            {loading ? (
              <div className="text-center py-8">
                <p className="text-gray-400">Loading users...</p>
              </div>
            ) : error ? (
              <div className="rounded-md bg-red-900/50 p-4">
                <p className="text-sm text-red-200">{error}</p>
              </div>
            ) : users.length === 0 ? (
              <div className="text-center py-8">
                <p className="text-gray-400">No users found</p>
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-700">
                  <thead>
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                        Username
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                        Email
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                        Roles
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                        Created
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-700">
                    {users.map((u) => (
                      <tr key={u.id} className="hover:bg-gray-700/50 transition-colors">
                        <td className="px-4 py-4 whitespace-nowrap text-sm font-medium text-white">
                          {u.username}
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-400">
                          {u.email}
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap text-sm">
                          <div className="flex flex-wrap gap-1">
                            {u.roles?.map((role) => (
                              <span
                                key={role}
                                className="px-2 py-1 bg-blue-600/20 text-blue-400 rounded text-xs font-medium"
                              >
                                {role}
                              </span>
                            ))}
                          </div>
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-400">
                          {new Date(u.created_at).toLocaleDateString()}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
