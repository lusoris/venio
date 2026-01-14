"use client";

import { useAuth } from "@/contexts/AuthContext";
import { redirect } from "next/navigation";
import { useState, useEffect } from "react";
import UserManagementTable from "@/components/admin/UserManagementTable";
import UserFormModal from "@/components/admin/UserFormModal";

interface User {
  id: string;
  email: string;
  name: string;
  roles: string[];
  created_at: string;
  verified: boolean;
}

export default function UsersPage() {
  const { user, isAuthenticated, token } = useAuth();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);

  if (!isAuthenticated || !user) {
    redirect("/login");
  }

  if (user.roles && !user.roles.includes("admin")) {
    redirect("/dashboard");
  }

  useEffect(() => {
    fetchUsers();
  }, [token]);

  const fetchUsers = async () => {
    try {
      setLoading(true);
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/admin/users`,
        {
          headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to fetch users: ${response.statusText}`);
      }

      const data = await response.json();
      setUsers(data.users || []);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch users");
    } finally {
      setLoading(false);
    }
  };

  const handleUserUpdated = () => {
    fetchUsers();
    setShowForm(false);
  };

  return (
    <div>
      <div className="mb-6 flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold text-white">User Management</h1>
          <p className="text-gray-400 mt-2">
            Manage users, roles, and permissions
          </p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
        >
          Add User
        </button>
      </div>

      {error && (
        <div className="bg-red-900 border border-red-700 text-red-100 px-4 py-3 rounded-lg mb-4">
          {error}
        </div>
      )}

      {loading ? (
        <div className="bg-gray-800 rounded-lg p-8 text-center border border-gray-700">
          <p className="text-gray-400">Loading users...</p>
        </div>
      ) : (
        <UserManagementTable
          users={users}
          onRefresh={fetchUsers}
          onUserUpdated={handleUserUpdated}
        />
      )}

      {showForm && (
        <UserFormModal
          onClose={() => setShowForm(false)}
          onSuccess={handleUserUpdated}
        />
      )}
    </div>
  );
}
