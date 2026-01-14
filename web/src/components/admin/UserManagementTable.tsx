"use client";

import { useAuth } from "@/contexts/AuthContext";
import { useState } from "react";

interface User {
  id: string;
  email: string;
  name: string;
  roles: string[];
  created_at: string;
  verified: boolean;
}

interface Props {
  users: User[];
  onRefresh: () => void;
  onUserUpdated: () => void;
}

export default function UserManagementTable({
  users,
  onRefresh,
  onUserUpdated,
}: Props) {
  const { token } = useAuth();
  const [loading, setLoading] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);

  const handleDelete = async (userId: string) => {
    if (!confirm("Are you sure you want to delete this user?")) return;

    try {
      setLoading(true);
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/admin/users/${userId}`,
        {
          method: "DELETE",
          headers: {
            "Authorization": `Bearer ${token}`,
          },
        }
      );

      if (response.ok) {
        onUserUpdated();
      } else {
        alert("Failed to delete user");
      }
    } catch (error) {
      console.error("Delete error:", error);
      alert("Error deleting user");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-gray-800 border border-gray-700 rounded-lg overflow-hidden">
      <div className="grid grid-cols-1 md:grid-cols-6 gap-4 p-4 bg-gray-700 font-semibold text-white text-sm">
        <div>Email</div>
        <div>Name</div>
        <div>Roles</div>
        <div>Status</div>
        <div>Created</div>
        <div>Actions</div>
      </div>
      <div className="divide-y divide-gray-700">
        {users.map((user) => (
          <div
            key={user.id}
            className="grid grid-cols-1 md:grid-cols-6 gap-4 p-4 hover:bg-gray-700 transition items-center"
          >
            <div className="text-white">{user.email}</div>
            <div className="text-gray-300">{user.name || "-"}</div>
            <div className="flex gap-1">
              {user.roles && user.roles.length > 0 ? (
                user.roles.map((role) => (
                  <span
                    key={role}
                    className="bg-blue-900 text-blue-200 px-2 py-1 rounded text-xs"
                  >
                    {role}
                  </span>
                ))
              ) : (
                <span className="text-gray-500">No roles</span>
              )}
            </div>
            <div>
              <span
                className={`px-2 py-1 rounded text-xs ${
                  user.verified
                    ? "bg-green-900 text-green-200"
                    : "bg-yellow-900 text-yellow-200"
                }`}
              >
                {user.verified ? "Verified" : "Pending"}
              </span>
            </div>
            <div className="text-gray-400 text-sm">
              {new Date(user.created_at).toLocaleDateString()}
            </div>
            <div className="flex gap-2">
              <button
                onClick={() => setEditingId(user.id)}
                className="text-blue-400 hover:text-blue-300 text-sm"
              >
                Edit
              </button>
              <button
                onClick={() => handleDelete(user.id)}
                disabled={loading}
                className="text-red-400 hover:text-red-300 text-sm disabled:opacity-50"
              >
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>
      {users.length === 0 && (
        <div className="p-8 text-center text-gray-400">No users found</div>
      )}
    </div>
  );
}
