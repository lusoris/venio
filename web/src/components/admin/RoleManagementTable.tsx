"use client";

import { useAuth } from "@/contexts/AuthContext";
import { useState } from "react";

interface Role {
  id: string;
  name: string;
  description: string;
  permissions: string[];
  user_count: number;
  created_at: string;
}

interface Props {
  roles: Role[];
  onRefresh: () => void;
  onRoleUpdated: () => void;
}

export default function RoleManagementTable({
  roles,
  onRefresh,
  onRoleUpdated,
}: Props) {
  const { token } = useAuth();
  const [loading, setLoading] = useState(false);

  const handleDelete = async (roleId: string) => {
    if (!confirm("Are you sure? This action cannot be undone.")) return;

    try {
      setLoading(true);
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/admin/roles/${roleId}`,
        {
          method: "DELETE",
          headers: {
            "Authorization": `Bearer ${token}`,
          },
        }
      );

      if (response.ok) {
        onRoleUpdated();
      } else {
        alert("Failed to delete role");
      }
    } catch (error) {
      console.error("Delete error:", error);
      alert("Error deleting role");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-gray-800 border border-gray-700 rounded-lg overflow-hidden">
      <div className="grid grid-cols-1 md:grid-cols-5 gap-4 p-4 bg-gray-700 font-semibold text-white text-sm">
        <div>Name</div>
        <div>Description</div>
        <div>Permissions</div>
        <div>Users</div>
        <div>Actions</div>
      </div>
      <div className="divide-y divide-gray-700">
        {roles.map((role) => (
          <div
            key={role.id}
            className="grid grid-cols-1 md:grid-cols-5 gap-4 p-4 hover:bg-gray-700 transition items-center"
          >
            <div>
              <span className="font-medium text-white">{role.name}</span>
            </div>
            <div className="text-gray-300 text-sm">{role.description}</div>
            <div>
              <span className="bg-purple-900 text-purple-200 px-2 py-1 rounded text-xs">
                {role.permissions.length} permissions
              </span>
            </div>
            <div>
              <span className="text-gray-300">{role.user_count} users</span>
            </div>
            <div className="flex gap-2">
              <button className="text-blue-400 hover:text-blue-300 text-sm">
                Edit
              </button>
              {role.user_count === 0 && (
                <button
                  onClick={() => handleDelete(role.id)}
                  disabled={loading}
                  className="text-red-400 hover:text-red-300 text-sm disabled:opacity-50"
                >
                  Delete
                </button>
              )}
            </div>
          </div>
        ))}
      </div>
      {roles.length === 0 && (
        <div className="p-8 text-center text-gray-400">No roles found</div>
      )}
    </div>
  );
}
