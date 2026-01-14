"use client";

import { useAuth } from "@/contexts/AuthContext";
import { useState } from "react";

interface UserRole {
  id: string;
  user_id: string;
  user_email: string;
  role_id: string;
  role_name: string;
  assigned_at: string;
  assigned_by: string;
}

interface Props {
  assignments: UserRole[];
  onRefresh: () => void;
}

export default function AssignmentManagementTable({ assignments, onRefresh }: Props) {
  const { token } = useAuth();
  const [loading, setLoading] = useState(false);

  const handleRevoke = async (assignmentId: string) => {
    if (!confirm("Revoke this role assignment?")) return;

    try {
      setLoading(true);
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/admin/user-roles/${assignmentId}`,
        {
          method: "DELETE",
          headers: {
            "Authorization": `Bearer ${token}`,
          },
        }
      );

      if (response.ok) {
        onRefresh();
      } else {
        alert("Failed to revoke assignment");
      }
    } catch (error) {
      console.error("Revoke error:", error);
      alert("Error revoking assignment");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-gray-800 border border-gray-700 rounded-lg overflow-hidden">
      <div className="grid grid-cols-1 md:grid-cols-6 gap-4 p-4 bg-gray-700 font-semibold text-white text-sm">
        <div>User</div>
        <div>Role</div>
        <div>Assigned By</div>
        <div>Assigned At</div>
        <div>Status</div>
        <div>Actions</div>
      </div>
      <div className="divide-y divide-gray-700">
        {assignments.map((assignment) => (
          <div
            key={assignment.id}
            className="grid grid-cols-1 md:grid-cols-6 gap-4 p-4 hover:bg-gray-700 transition items-center"
          >
            <div className="text-white">{assignment.user_email}</div>
            <div>
              <span className="bg-blue-900 text-blue-200 px-2 py-1 rounded text-xs">
                {assignment.role_name}
              </span>
            </div>
            <div className="text-gray-300">{assignment.assigned_by}</div>
            <div className="text-gray-400 text-sm">
              {new Date(assignment.assigned_at).toLocaleDateString()}
            </div>
            <div>
              <span className="bg-green-900 text-green-200 px-2 py-1 rounded text-xs">
                Active
              </span>
            </div>
            <div>
              <button
                onClick={() => handleRevoke(assignment.id)}
                disabled={loading}
                className="text-red-400 hover:text-red-300 text-sm disabled:opacity-50"
              >
                Revoke
              </button>
            </div>
          </div>
        ))}
      </div>
      {assignments.length === 0 && (
        <div className="p-8 text-center text-gray-400">
          No assignments found
        </div>
      )}
    </div>
  );
}
