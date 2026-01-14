"use client";

import { useAuth } from "@/contexts/AuthContext";
import { useState, useEffect } from "react";

interface Props {
  onClose: () => void;
  onSuccess: () => void;
}

export default function UserFormModal({ onClose, onSuccess }: Props) {
  const { token } = useAuth();
  const [loading, setLoading] = useState(false);
  const [roles, setRoles] = useState<any[]>([]);
  const [formData, setFormData] = useState({
    email: "",
    name: "",
    password: "",
    roles: [] as string[],
  });

  useEffect(() => {
    fetchRoles();
  }, []);

  const fetchRoles = async () => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/admin/roles`,
        {
          headers: {
            "Authorization": `Bearer ${token}`,
          },
        }
      );
      const data = await response.json();
      setRoles(data.roles || []);
    } catch (error) {
      console.error("Error fetching roles:", error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/admin/users`,
        {
          method: "POST",
          headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify(formData),
        }
      );

      if (response.ok) {
        onSuccess();
      } else {
        alert("Failed to create user");
      }
    } catch (error) {
      console.error("Error creating user:", error);
      alert("Error creating user");
    } finally {
      setLoading(false);
    }
  };

  const toggleRole = (roleId: string) => {
    setFormData((prev) => ({
      ...prev,
      roles: prev.roles.includes(roleId)
        ? prev.roles.filter((r) => r !== roleId)
        : [...prev.roles, roleId],
    }));
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-gray-800 rounded-lg p-6 w-full max-w-md">
        <h2 className="text-xl font-bold text-white mb-4">Create New User</h2>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">
              Email
            </label>
            <input
              type="email"
              required
              value={formData.email}
              onChange={(e) =>
                setFormData((prev) => ({ ...prev, email: e.target.value }))
              }
              className="w-full bg-gray-700 text-white px-3 py-2 rounded-lg border border-gray-600 focus:border-blue-500 focus:outline-none"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">
              Name
            </label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) =>
                setFormData((prev) => ({ ...prev, name: e.target.value }))
              }
              className="w-full bg-gray-700 text-white px-3 py-2 rounded-lg border border-gray-600 focus:border-blue-500 focus:outline-none"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">
              Password
            </label>
            <input
              type="password"
              required
              value={formData.password}
              onChange={(e) =>
                setFormData((prev) => ({ ...prev, password: e.target.value }))
              }
              className="w-full bg-gray-700 text-white px-3 py-2 rounded-lg border border-gray-600 focus:border-blue-500 focus:outline-none"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Assign Roles
            </label>
            <div className="space-y-2">
              {roles.map((role) => (
                <label key={role.id} className="flex items-center">
                  <input
                    type="checkbox"
                    checked={formData.roles.includes(role.id)}
                    onChange={() => toggleRole(role.id)}
                    className="w-4 h-4 rounded"
                  />
                  <span className="ml-2 text-gray-300">{role.name}</span>
                </label>
              ))}
            </div>
          </div>

          <div className="flex gap-2 pt-4">
            <button
              type="submit"
              disabled={loading}
              className="flex-1 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg disabled:opacity-50"
            >
              Create User
            </button>
            <button
              type="button"
              onClick={onClose}
              className="flex-1 bg-gray-700 hover:bg-gray-600 text-white px-4 py-2 rounded-lg"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
