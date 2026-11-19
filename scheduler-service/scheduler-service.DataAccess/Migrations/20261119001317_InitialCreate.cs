using System;
using Microsoft.EntityFrameworkCore.Migrations;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;

#nullable disable

namespace scheduler_service.DataAccess.Migrations
{
    /// <inheritdoc />
    public partial class InitialCreate : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "reserving",
                columns: table => new
                {
                    id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    name = table.Column<string>(type: "text", nullable: false),
                    Age = table.Column<int>(type: "integer", nullable: false),
                    email = table.Column<string>(type: "text", nullable: false),
                    phone_number = table.Column<int>(type: "integer", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_reserving", x => x.id);
                });

            migrationBuilder.CreateTable(
                name: "timeslot",
                columns: table => new
                {
                    id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    apartment_id = table.Column<long>(type: "bigint", nullable: false),
                    startdatetime = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    enddatetime = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_timeslot", x => x.id);
                });

            migrationBuilder.CreateTable(
                name: "reservation",
                columns: table => new
                {
                    id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    timeslot_id = table.Column<long>(type: "bigint", nullable: false),
                    reservation_status = table.Column<int>(type: "integer", nullable: false),
                    reserving_id = table.Column<long>(type: "bigint", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_reservation", x => x.id);
                    table.ForeignKey(
                        name: "FK_reservation_reserving_reserving_id",
                        column: x => x.reserving_id,
                        principalTable: "reserving",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "FK_reservation_timeslot_timeslot_id",
                        column: x => x.timeslot_id,
                        principalTable: "timeslot",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateIndex(
                name: "IX_reservation_reserving_id",
                table: "reservation",
                column: "reserving_id",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "IX_reservation_timeslot_id",
                table: "reservation",
                column: "timeslot_id");

            migrationBuilder.CreateIndex(
                name: "IX_timeslot_apartment_id",
                table: "timeslot",
                column: "apartment_id");

            migrationBuilder.CreateIndex(
                name: "IX_timeslot_startdatetime_enddatetime",
                table: "timeslot",
                columns: new[] { "startdatetime", "enddatetime" });
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "reservation");

            migrationBuilder.DropTable(
                name: "reserving");

            migrationBuilder.DropTable(
                name: "timeslot");
        }
    }
}
